package message

import (
	"time"
)

type Message struct {
	ID        int64         `json:"id"`
	To        string        `json:"to"`
	Content   string        `json:"content"`
	Status    MessageStatus `json:"status"`
	SentAt    *time.Time    `json:"sent_at,omitempty"`
	MessageID *string       `json:"message_id,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
}

type MessageStatus string

const (
	StatusPending MessageStatus = "pending"
	StatusSent    MessageStatus = "sent"
	StatusFailed  MessageStatus = "failed"
	MaxContentLen int           = 160
)
