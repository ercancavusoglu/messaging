package webhook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ercancavusoglu/messaging/internal/domain/message"
	"github.com/ercancavusoglu/messaging/pkg/logger"
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

type webhookRequest struct {
	To      string `json:"to"`
	Content string `json:"content"`
}

func (c *Client) SendMessage(to, content string) (*message.WebhookResponse, error) {
	payload := webhookRequest{
		To:      to,
		Content: content,
	}

	logger.Info("Sending message to webhook", "to", to, "content", content)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Error marshaling request", "error", err)
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Error creating request", "error", err)
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error sending request", "error", err)
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response message.WebhookResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logger.Error("Error decoding response", "error", err)
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	logger.Info("Message sent successfully", "messageId", response.MessageID)

	return &response, nil
}
