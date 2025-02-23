package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ercancavusoglu/messaging/internal/domain"
	"io"
	"log"
	"net/http"
)

type Client struct {
	url    string
	apiKey string
}

func NewClient(url, apiKey string) *Client {
	return &Client{
		url:    url,
		apiKey: apiKey,
	}
}

func (c *Client) SendMessage(to, content string) (*domain.WebhookResponse, error) {
	log.Printf("[Webhook] Sending message [to: %s, content: %s]", to, content)

	payload := map[string]string{
		"to":      to,
		"content": content,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	log.Printf("[Webhook] Request payload: %s", string(jsonPayload))

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

	log.Printf("[Webhook] Sending request to: %s", c.url)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	log.Printf("[Webhook] Response status: %d", resp.StatusCode)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}
	log.Printf("[Webhook] Response body: %s", string(bodyBytes))

	var response domain.WebhookResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	log.Printf("[Webhook] Decoded response: %+v", response)
	return &response, nil
}
