package webhook

import (
	"fmt"

	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/ercancavusoglu/messaging/internal/ports"
)

type RetryableWebhookClient struct {
	clients    []ports.WebhookClient
	maxRetries int
}

func NewRetryableWebhookClient(clients []ports.WebhookClient, maxRetries int) *RetryableWebhookClient {
	return &RetryableWebhookClient{
		clients:    clients,
		maxRetries: maxRetries,
	}
}

func (c *RetryableWebhookClient) SendMessage(to, content string) (*domain.WebhookResponse, error) {
	var lastErr error

	fmt.Println("Sending message to", to, "with content", content)
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		clientIndex := attempt % len(c.clients)
		client := c.clients[clientIndex]

		response, err := client.SendMessage(to, content)
		if err == nil {
			fmt.Println("attempt", attempt+1, "succeeded with client", clientIndex+1)
			return response, nil
		}

		lastErr = fmt.Errorf("attempt %d failed with client %d: %v", attempt+1, clientIndex+1, err)
		fmt.Println(lastErr)
	}

	return nil, fmt.Errorf("all retry attempts failed after %d tries. Last error: %v", c.maxRetries, lastErr)
}
