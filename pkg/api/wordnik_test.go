package api

import (
	"net/http"
	"testing"
)

func TestCalculateRarityTier(t *testing.T) {
	tests := []struct {
		corpusCount int64
		zipf        float64
		expected    string
	}{
		{30, 1.8, "Rare"},
		{150, 2.5, "Obscure"},
		{350, 3.0, "Elegant"},
		{1000, 3.8, "Uncommon"},
	}

	for _, tt := range tests {
		got := CalculateRarityTier(tt.corpusCount, tt.zipf)
		if got != tt.expected {
			t.Errorf("For count %d, zipf %.1f expected %s, got %s", tt.corpusCount, tt.zipf, tt.expected, got)
		}
	}
}

func TestParseRateLimitHeaders(t *testing.T) {
	client := NewClient("test_key")
	header := http.Header{}
	header.Set("x-ratelimit-remaining-hour", "450")
	header.Set("x-ratelimit-remaining-minute", "45")
	header.Set("x-ratelimit-limit-minute", "50")
	header.Set("x-ratelimit-limit-hour", "500")

	client.parseRateLimitHeaders(header)
	info := client.GetRateLimitInfo()

	if info.RemainingHour != 450 {
		t.Errorf("Expected remaining hour 450, got %d", info.RemainingHour)
	}
	if info.RemainingMinute != 45 {
		t.Errorf("Expected remaining minute 45, got %d", info.RemainingMinute)
	}
}
