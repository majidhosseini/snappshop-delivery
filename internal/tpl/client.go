package tpl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	CreateShipment(ctx context.Context, orderID string) error
}

type client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) Client {
	return &client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // Prevent long-running requests
		},
	}
}

// CreateShipment sends order data to 3PL API
func (c *client) CreateShipment(ctx context.Context, orderID string) error {
	// Example request payload
	payload := map[string]interface{}{
		"order_id":  orderID,
		"timestamp": time.Now().UTC(),
	}

	jsonBody, _ := json.Marshal(payload)

	req, _ := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/shipments", c.baseURL),
		bytes.NewReader(jsonBody),
	)

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("3PL API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResp struct {
			Error string `json:"error"`
		}
		json.NewDecoder(resp.Body).Decode(&errorResp)
		return fmt.Errorf("3PL API error: %s (status %d)", errorResp.Error, resp.StatusCode)
	}

	return nil
}
