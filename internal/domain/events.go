package domain

import (
	"fmt"
	"strconv"
)

const (
	EventMessageSent   = "message.sent"
	EventMessageFailed = "message.failed"
	EventMessageQueued = "message.queued"
)

type MessageSentEvent struct {
	BaseEvent
	Message   *Message `json:"message"`
	MessageID string   `json:"message_id"`
}

func NewMessageSentEvent(message *Message, messageID string) MessageSentEvent {
	return MessageSentEvent{
		BaseEvent: NewBaseEvent(EventMessageSent, strconv.FormatInt(message.ID, 10)),
		Message:   message,
		MessageID: messageID,
	}
}

type MessageFailedEvent struct {
	BaseEvent
	Message *Message `json:"message"`
	Error   string   `json:"error"`
}

func NewMessageFailedEvent(message *Message, err error) MessageFailedEvent {
	return MessageFailedEvent{
		BaseEvent: NewBaseEvent(EventMessageFailed, strconv.FormatInt(message.ID, 10)),
		Message:   message,
		Error:     fmt.Sprintf("%v", err),
	}
}

type MessageQueuedEvent struct {
	BaseEvent
	Message *Message `json:"message"`
}

func NewMessageQueuedEvent(message *Message) MessageQueuedEvent {
	return MessageQueuedEvent{
		BaseEvent: NewBaseEvent(EventMessageQueued, strconv.FormatInt(message.ID, 10)),
		Message:   message,
	}
}
