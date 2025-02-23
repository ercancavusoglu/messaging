package ports

import "github.com/ercancavusoglu/messaging/internal/domain"

type WebhookClient interface {
	SendMessage(to, content string) (*domain.WebhookResponse, error)
}
