package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents the application configuration
type Config struct {
	ServerURL      string        `json:"serverUrl"`
	Username       string        `json:"username"`
	Password       string        `json:"password,omitempty"` // Consider more secure storage
	SyncFrequency  time.Duration `json:"syncFrequency"`
	WatchDirs      []string      `json:"watchDirs"`
	IgnorePatterns []string      `json:"ignorePatterns"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Errorf("failed to get home directory: %w", err))
	}

	return &Config{
		ServerURL:     "http://localhost:8080",
		SyncFrequency: 5 * time.Minute,
		WatchDirs:     []string{filepath.Join(homeDir, "homecloud")},
		IgnorePatterns: []string{
			".DS_Store",
			"Thumbs.db",
			"*.tmp",
		},
	}
}

// LoadConfig loads configuration from the specified path
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Create a default config if none exists
			config := DefaultConfig()
			if err := SaveConfig(path, config); err != nil {
				return nil, fmt.Errorf("failed to create default config: %w", err)
			}
			return config, nil
		}
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the specified path
func SaveConfig(path string, config *Config) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}