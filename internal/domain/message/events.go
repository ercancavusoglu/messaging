package message

import (
	"fmt"
	"strconv"

	"github.com/ercancavusoglu/messaging/internal/domain/event"
)

const (
	EventMessageCreated  = "message.created"
	EventMessageSent     = "message.sent"
	EventMessageFailed   = "message.failed"
	EventMessageQueued   = "message.queued"
	EventMessageDequeued = "message.dequeued"
)

type MessageCreatedEvent struct {
	event.BaseEvent
	Message *Message `json:"message"`
}

func NewMessageCreatedEvent(message *Message) MessageCreatedEvent {
	return MessageCreatedEvent{
		BaseEvent: event.NewBaseEvent(EventMessageCreated, strconv.FormatInt(message.ID, 10)),
		Message:   message,
	}
}

type MessageSentEvent struct {
	event.BaseEvent
	Message   *Message `json:"message"`
	MessageID string   `json:"message_id"`
}

func NewMessageSentEvent(message *Message, messageID string) MessageSentEvent {
	return MessageSentEvent{
		BaseEvent: event.NewBaseEvent(EventMessageSent, strconv.FormatInt(message.ID, 10)),
		Message:   message,
		MessageID: messageID,
	}
}

type MessageFailedEvent struct {
	event.BaseEvent
	Message *Message `json:"message"`
	Error   string   `json:"error"`
}

func NewMessageFailedEvent(message *Message, err error) MessageFailedEvent {
	return MessageFailedEvent{
		BaseEvent: event.NewBaseEvent(EventMessageFailed, strconv.FormatInt(message.ID, 10)),
		Message:   message,
		Error:     fmt.Sprintf("%v", err),
	}
}

type MessageQueuedEvent struct {
	event.BaseEvent
	Message *Message `json:"message"`
}

func NewMessageQueuedEvent(message *Message) MessageQueuedEvent {
	return MessageQueuedEvent{
		BaseEvent: event.NewBaseEvent(EventMessageQueued, strconv.FormatInt(message.ID, 10)),
		Message:   message,
	}
}

type MessageDequeuedEvent struct {
	event.BaseEvent
	Message *Message `json:"message"`
}

func NewMessageDequeuedEvent(message *Message) MessageDequeuedEvent {
	return MessageDequeuedEvent{
		BaseEvent: event.NewBaseEvent(EventMessageDequeued, strconv.FormatInt(message.ID, 10)),
		Message:   message,
	}
}
