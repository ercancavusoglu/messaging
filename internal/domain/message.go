package domain

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
