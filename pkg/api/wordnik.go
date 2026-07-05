package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/wordy-tui/wordy/pkg/seed"
)

const (
	BaseURL = "https://api.wordnik.com/v4"
)

type RateLimitInfo struct {
	RemainingHour   int       `json:"remaining_hour"`
	RemainingMinute int       `json:"remaining_minute"`
	LimitMinute     int       `json:"limit_minute"`
	LimitHour       int       `json:"limit_hour"`
	IsRateLimited   bool      `json:"is_rate_limited"`
	ResetTime       time.Time `json:"reset_time"`
}

type Client struct {
	apiKey     string
	httpClient *http.Client
	cache      *DiskCache
	mu         sync.RWMutex
	rateLimit  RateLimitInfo
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		cache: NewDiskCache(),
		rateLimit: RateLimitInfo{
			RemainingHour:   500,
			RemainingMinute: 50,
			LimitMinute:     50,
			LimitHour:       500,
		},
	}
}

func (c *Client) SetAPIKey(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.apiKey = key
}

func (c *Client) GetRateLimitInfo() RateLimitInfo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rateLimit
}

func (c *Client) parseRateLimitHeaders(headers http.Header) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if val := headers.Get("x-ratelimit-remaining-hour"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			c.rateLimit.RemainingHour = v
		}
	}
	if val := headers.Get("x-ratelimit-remaining-minute"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			c.rateLimit.RemainingMinute = v
		}
	}
	if val := headers.Get("x-ratelimit-limit-minute"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			c.rateLimit.LimitMinute = v
		}
	}
	if val := headers.Get("x-ratelimit-limit-hour"); val != "" {
		if v, err := strconv.Atoi(val); err == nil {
			c.rateLimit.LimitHour = v
		}
	}
}

// FetchWordDetails retrieves comprehensive word details from cache, seed, or live Wordnik API.
func (c *Client) FetchWordDetails(word string) (seed.WordDetails, error) {
	wordLower := strings.TrimSpace(strings.ToLower(word))
	if wordLower == "" {
		return seed.WordDetails{}, fmt.Errorf("empty word")
	}

	// 1. Check local disk cache
	if cached, ok := c.cache.Get(wordLower); ok {
		return cached, nil
	}

	// 2. Check seed dataset
	for _, s := range seed.GetSeedWords() {
		if strings.ToLower(s.Word) == wordLower {
			c.cache.Set(wordLower, s)
			return s, nil
		}
	}

	// 3. If no API key, return offline placeholder error
	if c.apiKey == "" {
		return seed.WordDetails{}, fmt.Errorf("no API key set (operating in offline seed mode)")
	}

	// 4. Query Wordnik endpoints concurrently / sequentially
	details, err := c.fetchLiveWordnikData(wordLower)
	if err != nil {
		return seed.WordDetails{}, err
	}

	// Save to disk cache
	c.cache.Set(wordLower, details)
	return details, nil
}

func (c *Client) fetchLiveWordnikData(word string) (seed.WordDetails, error) {
	details := seed.WordDetails{
		Word:         word,
		RelatedWords: make(map[string][]string),
		RarityTier:   "Uncommon",
	}

	// Definitions
	defs, attrText, err := c.getDefinitions(word)
	if err == nil && len(defs) > 0 {
		details.Definitions = defs
		details.PartOfSpeech = defs[0].PartOfSpeech
		details.AttributionText = attrText
	} else if err != nil && strings.Contains(err.Error(), "429") {
		return details, err
	}

	// Frequency
	corpusCount, zipf, tier := c.getFrequency(word)
	details.CorpusCount = corpusCount
	details.ZipfScore = zipf
	details.RarityTier = tier

	// Related words
	rel := c.getRelatedWords(word)
	if len(rel) > 0 {
		details.RelatedWords = rel
	}

	// Pronunciations
	phonetic := c.getPronunciation(word)
	if phonetic != "" {
		details.Phonetic = phonetic
	}

	// Examples
	examples := c.getExamples(word)
	if len(examples) > 0 {
		details.Examples = examples
	}

	return details, nil
}

func (c *Client) getDefinitions(word string) ([]seed.Definition, string, error) {
	u := fmt.Sprintf("%s/word.json/%s/definitions?limit=5&useCanonical=true&api_key=%s",
		BaseURL, url.PathEscape(word), c.apiKey)

	resp, err := c.httpClient.Get(u)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	c.parseRateLimitHeaders(resp.Header)

	if resp.StatusCode == 429 {
		c.mu.Lock()
		c.rateLimit.IsRateLimited = true
		c.rateLimit.ResetTime = time.Now().Add(1 * time.Minute)
		c.mu.Unlock()
		return nil, "", fmt.Errorf("rate limited (429)")
	}

	if resp.StatusCode != 200 {
		return nil, "", fmt.Errorf("status code %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var rawDefs []struct {
		Text            string `json:"text"`
		PartOfSpeech    string `json:"partOfSpeech"`
		AttributionText string `json:"attributionText"`
	}

	if err := json.Unmarshal(body, &rawDefs); err != nil {
		return nil, "", err
	}

	defs := make([]seed.Definition, 0)
	attrText := ""
	for _, rd := range rawDefs {
		if rd.Text != "" {
			defs = append(defs, seed.Definition{
				Text:            rd.Text,
				PartOfSpeech:    rd.PartOfSpeech,
				AttributionText: rd.AttributionText,
			})
			if attrText == "" && rd.AttributionText != "" {
				attrText = rd.AttributionText
			}
		}
	}
	return defs, attrText, nil
}

func (c *Client) getFrequency(word string) (int64, float64, string) {
	u := fmt.Sprintf("%s/word.json/%s/frequency?useCanonical=true&api_key=%s",
		BaseURL, url.PathEscape(word), c.apiKey)

	resp, err := c.httpClient.Get(u)
	if err != nil || resp.StatusCode != 200 {
		return 300, 3.0, "Uncommon"
	}
	defer resp.Body.Close()

	c.parseRateLimitHeaders(resp.Header)

	var raw struct {
		TotalCount int64 `json:"totalCount"`
		Frequency  []struct {
			Count int64 `json:"count"`
			Year  int   `json:"year"`
		} `json:"frequency"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &raw); err != nil {
		return 300, 3.0, "Uncommon"
	}

	total := raw.TotalCount
	if total == 0 {
		var sum int64
		for _, f := range raw.Frequency {
			sum += f.Count
		}
		total = sum
	}

	// Calculate approximate Zipf score
	zipf := 3.0
	if total > 0 {
		zipf = math.Round((math.Log10(float64(total))+1.5)*10) / 10
	}

	tier := CalculateRarityTier(total, zipf)
	return total, zipf, tier
}

func (c *Client) getRelatedWords(word string) map[string][]string {
	u := fmt.Sprintf("%s/word.json/%s/relatedWords?useCanonical=true&limitPerRelationshipType=8&api_key=%s",
		BaseURL, url.PathEscape(word), c.apiKey)

	resp, err := c.httpClient.Get(u)
	if err != nil || resp.StatusCode != 200 {
		return nil
	}
	defer resp.Body.Close()

	c.parseRateLimitHeaders(resp.Header)

	var raw []struct {
		RelationshipType string   `json:"relationshipType"`
		Words            []string `json:"words"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}

	res := make(map[string][]string)
	for _, item := range raw {
		if len(item.Words) > 0 {
			res[item.RelationshipType] = item.Words
		}
	}
	return res
}

func (c *Client) getPronunciation(word string) string {
	u := fmt.Sprintf("%s/word.json/%s/pronunciations?useCanonical=true&limit=1&api_key=%s",
		BaseURL, url.PathEscape(word), c.apiKey)

	resp, err := c.httpClient.Get(u)
	if err != nil || resp.StatusCode != 200 {
		return ""
	}
	defer resp.Body.Close()

	c.parseRateLimitHeaders(resp.Header)

	var raw []struct {
		Raw     string `json:"raw"`
		RawType string `json:"rawType"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &raw); err == nil && len(raw) > 0 {
		return raw[0].Raw
	}
	return ""
}

func (c *Client) getExamples(word string) []seed.Example {
	u := fmt.Sprintf("%s/word.json/%s/examples?limit=3&api_key=%s",
		BaseURL, url.PathEscape(word), c.apiKey)

	resp, err := c.httpClient.Get(u)
	if err != nil || resp.StatusCode != 200 {
		return nil
	}
	defer resp.Body.Close()

	c.parseRateLimitHeaders(resp.Header)

	var raw struct {
		Examples []struct {
			Text   string `json:"text"`
			Title  string `json:"title"`
			Author string `json:"author"`
		} `json:"examples"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil
	}

	examples := make([]seed.Example, 0)
	for _, ex := range raw.Examples {
		examples = append(examples, seed.Example{
			Text:   ex.Text,
			Title:  ex.Title,
			Author: ex.Author,
		})
	}
	return examples
}

func CalculateRarityTier(corpusCount int64, zipf float64) string {
	if corpusCount < 50 || zipf < 2.0 {
		return "Rare"
	} else if corpusCount < 200 || zipf < 2.7 {
		return "Obscure"
	} else if corpusCount < 500 || zipf < 3.2 {
		return "Elegant"
	}
	return "Uncommon"
}
