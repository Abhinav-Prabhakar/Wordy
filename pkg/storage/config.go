package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds user configuration for Wordy.
type Config struct {
	WordnikAPIKey    string `json:"wordnik_api_key"`
	MinCorpusCount   int    `json:"min_corpus_count"`
	MaxCorpusCount   int    `json:"max_corpus_count"`
	TargetRarityTier string `json:"target_rarity_tier"` // "all", "uncommon", "rare", "obscure"
	Theme            string `json:"theme"`              // "bubbletea", "dracula", "nord"
}

// DefaultConfig returns reasonable default settings.
func DefaultConfig() Config {
	return Config{
		WordnikAPIKey:    os.Getenv("WORDNIK_API_KEY"),
		MinCorpusCount:   10,
		MaxCorpusCount:   10000,
		TargetRarityTier: "uncommon",
		Theme:            "bubbletea",
	}
}

// GetConfigPath returns the filepath to config.json.
func GetConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		homeDir, _ := os.UserHomeDir()
		configDir = filepath.Join(homeDir, ".config")
	}
	return filepath.Join(configDir, "wordy", "config.json")
}

// LoadConfig reads configuration from file or returns defaults.
func LoadConfig() Config {
	cfg := DefaultConfig()

	// Environment variable overrides file config if present
	if envKey := os.Getenv("WORDNIK_API_KEY"); envKey != "" {
		cfg.WordnikAPIKey = envKey
	}

	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}

	var loaded Config
	if err := json.Unmarshal(data, &loaded); err == nil {
		if envKey := os.Getenv("WORDNIK_API_KEY"); envKey != "" {
			loaded.WordnikAPIKey = envKey
		}
		return loaded
	}
	return cfg
}

// SaveConfig writes configuration to disk.
func SaveConfig(cfg Config) error {
	path := GetConfigPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
