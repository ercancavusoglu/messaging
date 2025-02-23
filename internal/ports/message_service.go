package ports

import "github.com/ercancavusoglu/messaging/internal/domain"

type MessageService interface {
	GetPendingMessages() ([]*domain.Message, error)
	GetSendedMessages() ([]*domain.Message, error)
	QueueMessage(msg *domain.Message) error
}
