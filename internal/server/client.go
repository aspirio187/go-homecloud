package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"homecloud/internal/models"
)

// Client handles communication with the remote server
type Client struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
}

// NewClient creates a new server client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Authenticate authenticates with the server and stores the token
func (c *Client) Authenticate(username, password string) error {
	authData := map[string]string{
		"username": username,
		"password": password,
	}

	data, err := json.Marshal(authData)
	if err != nil {
		return fmt.Errorf("failed to marshal auth data: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/auth/login", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed: status code %d", resp.StatusCode)
	}

	var result struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse authentication response: %w", err)
	}

	c.authToken = result.Token
	return nil
}

// UploadFile uploads a file to the server
func (c *Client) UploadFile(path string, content []byte, metadata map[string]string) error {
	if c.authToken == "" {
		return fmt.Errorf("not authenticated")
	}

	body := new(bytes.Buffer)
	// writer, err := http.NewRequest("POST", c.baseURL+"/api/files/upload", body)
	// TODO: Implement multipart form upload for file

	req, err := http.NewRequest("POST", c.baseURL+"/api/files/upload", body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("upload request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("upload failed: status code %d", resp.StatusCode)
	}

	return nil
}

// DownloadFile downloads a file from the server
func (c *Client) DownloadFile(path string) ([]byte, error) {
	if c.authToken == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	req, err := http.NewRequest("GET", c.baseURL+"/api/files/download?path="+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed: status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// GetFileMetadata retrieves file metadata from the server
func (c *Client) GetFileMetadata(path string) (map[string]string, error) {
	if c.authToken == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	req, err := http.NewRequest("GET", c.baseURL+"/api/files/metadata?path="+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("metadata request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("metadata request failed: status code %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse metadata response: %w", err)
	}

	return result, nil
}

// ListFiles lists files from the server
func (c *Client) ListFiles(path string) ([]*models.FileInfo, error) {
	if c.authToken == "" {
		return nil, fmt.Errorf("not authenticated")
	}

	req, err := http.NewRequest("GET", c.baseURL+"/api/files/list?path="+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list files request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("list files request failed: status code %d", resp.StatusCode)
	}

	var result []*models.FileInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse list files response: %w", err)
	}

	return result, nil
}

// DeleteFile deletes a file from the server
func (c *Client) DeleteFile(path string) error {
	if c.authToken == "" {
		return fmt.Errorf("not authenticated")
	}

	req, err := http.NewRequest("DELETE", c.baseURL+"/api/files?path="+path, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.authToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete failed: status code %d", resp.StatusCode)
	}

	return nil
}
