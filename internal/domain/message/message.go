package message

import "time"

type MessageStatus string

const (
	StatusPending MessageStatus = "pending"
	StatusSent    MessageStatus = "sent"
	StatusFailed  MessageStatus = "failed"
)

type Message struct {
	ID        int64         `json:"id"`
	To        string        `json:"to"`
	Content   string        `json:"content"`
	Status    MessageStatus `json:"status"`
	MessageID string        `json:"message_id"`
	CreatedAt time.Time     `json:"created_at"`
	SentAt    *time.Time    `json:"sent_at,omitempty"`
}

type WebhookResponse struct {
	MessageID string `json:"messageId"`
	Message   string `json:"message"`
}

type Repository interface {
	Create(message *Message) error
	GetByID(id int64) (*Message, error)
	GetPendingMessages(limit int) ([]*Message, error)
	UpdateStatus(id int64, status MessageStatus, messageID string) error
	GetByStatus(status MessageStatus) ([]*Message, error)
	List() ([]*Message, error)
}

type WebhookClient interface {
	SendMessage(to, content string) (*WebhookResponse, error)
}

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}
