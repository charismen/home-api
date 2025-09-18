package apiclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents an API client with configurable timeouts and retries
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	MaxRetries int
}

// NewClient creates a new API client
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		MaxRetries: 3,
	}
}

// FetchItems fetches items from the external API
// For this example, we'll use the PokéAPI to fetch Pokémon data
func (c *Client) FetchItems(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/pokemon?limit=%d", c.BaseURL, limit)
	
	var result struct {
		Results []map[string]interface{} `json:"results"`
	}
	
	err := c.getWithRetry(ctx, url, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch items: %w", err)
	}
	
	// For each result, fetch the detailed data
	var items []map[string]interface{}
	for _, item := range result.Results {
		url, ok := item["url"].(string)
		if !ok {
			continue
		}
		
		var detailedItem map[string]interface{}
		err := c.getWithRetry(ctx, url, &detailedItem)
		if err != nil {
			// Log error but continue with other items
			fmt.Printf("Error fetching detailed item data: %v\n", err)
			continue
		}
		
		// Add the detailed data
		items = append(items, detailedItem)
	}
	
	return items, nil
}

// getWithRetry performs a GET request with retry logic
func (c *Client) getWithRetry(ctx context.Context, url string, result interface{}) error {
	var lastErr error
	
	for attempt := 0; attempt <= c.MaxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s, ...
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			time.Sleep(backoff)
		}
		
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			lastErr = err
			continue
		}
		
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		
		defer resp.Body.Close()
		
		// Check if the response was successful
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				lastErr = err
				continue
			}
			
			err = json.Unmarshal(body, result)
			if err != nil {
				lastErr = err
				continue
			}
			
			return nil // Success
		}
		
		// Handle specific status codes
		switch resp.StatusCode {
		case http.StatusTooManyRequests, http.StatusServiceUnavailable:
			// These are retryable errors
			lastErr = fmt.Errorf("server returned %d, will retry", resp.StatusCode)
		default:
			// For other errors, read the response body for more details
			body, _ := io.ReadAll(resp.Body)
			lastErr = fmt.Errorf("server returned %d: %s", resp.StatusCode, string(body))
		}
	}
	
	return fmt.Errorf("after %d attempts: %w", c.MaxRetries+1, lastErr)
}