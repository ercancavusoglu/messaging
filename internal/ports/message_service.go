package ports

import "github.com/ercancavusoglu/messaging/internal/domain"

type MessageService interface {
	GetPendingMessages() ([]*domain.Message, error)
	GetSendedMessages() ([]*domain.Message, error)
	Publish(msg *domain.Message) error
}
