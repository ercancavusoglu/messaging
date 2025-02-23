package ports

import "github.com/ercancavusoglu/messaging/internal/domain"

type Repository interface {
	GetPendingMessages() ([]*domain.Message, error)
	UpdateStatus(id int64, status domain.MessageStatus, messageID string) error
	GetByStatus(status domain.MessageStatus) ([]*domain.Message, error)
}
