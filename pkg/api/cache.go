package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wordy-tui/wordy/pkg/seed"
)

type CachedItem struct {
	Timestamp time.Time        `json:"timestamp"`
	Data      seed.WordDetails `json:"data"`
}

type DiskCache struct {
	mu   sync.RWMutex
	path string
	Items map[string]CachedItem `json:"items"`
}

func GetCachePath() string {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		homeDir, _ := os.UserHomeDir()
		cacheDir = filepath.Join(homeDir, ".cache")
	}
	return filepath.Join(cacheDir, "wordy", "cache.json")
}

func NewDiskCache() *DiskCache {
	c := &DiskCache{
		path:  GetCachePath(),
		Items: make(map[string]CachedItem),
	}
	_ = c.Load()
	return c
}

func (c *DiskCache) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := os.ReadFile(c.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &c.Items)
}

func (c *DiskCache) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	dir := filepath.Dir(c.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c.Items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.path, data, 0644)
}

func (c *DiskCache) Get(word string) (seed.WordDetails, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.Items[word]
	if !exists {
		return seed.WordDetails{}, false
	}
	// Expire after 30 days
	if time.Since(item.Timestamp) > 30*24*time.Hour {
		return seed.WordDetails{}, false
	}
	return item.Data, true
}

func (c *DiskCache) Set(word string, data seed.WordDetails) {
	c.mu.Lock()
	c.Items[word] = CachedItem{
		Timestamp: time.Now(),
		Data:      data,
	}
	c.mu.Unlock()
	_ = c.Save()
}
