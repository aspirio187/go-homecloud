package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// AuthInfo stores authentication information
type AuthInfo struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	Username  string    `json:"username"`
}

// LoadAuth loads authentication information from disk
func LoadAuth(appDataPath string) (*AuthInfo, error) {
	authPath := filepath.Join(appDataPath, "auth.json")
	
	data, err := os.ReadFile(authPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read auth file: %w", err)
	}
	
	var authInfo AuthInfo
	if err := json.Unmarshal(data, &authInfo); err != nil {
		return nil, fmt.Errorf("failed to parse auth data: %w", err)
	}
	
	// Check if token is expired
	if authInfo.ExpiresAt.Before(time.Now()) {
		return nil, nil
	}
	
	return &authInfo, nil
}

// SaveAuth saves authentication information to disk
func SaveAuth(appDataPath string, authInfo *AuthInfo) error {
	authPath := filepath.Join(appDataPath, "auth.json")
	
	// Ensure directory exists
	dir := filepath.Dir(authPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create auth directory: %w", err)
	}
	
	data, err := json.MarshalIndent(authInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal auth info: %w", err)
	}
	
	if err := os.WriteFile(authPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write auth file: %w", err)
	}
	
	return nil
}

// ClearAuth removes saved authentication information
func ClearAuth(appDataPath string) error {
	authPath := filepath.Join(appDataPath, "auth.json")
	
	// Remove auth file if it exists
	if _, err := os.Stat(authPath); err == nil {
		if err := os.Remove(authPath); err != nil {
			return fmt.Errorf("failed to remove auth file: %w", err)
		}
	}
	
	return nil
}