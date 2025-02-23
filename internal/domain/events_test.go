package domain

import (
	"errors"
	"testing"
	"time"
)

func createTestMessage() *Message {
	return &Message{
		ID:        123,
		To:        "+905551234567",
		Content:   "Test message",
		Status:    StatusPending,
		MessageID: "",
		Provider:  "webhook",
		CreatedAt: time.Now(),
		SentAt:    nil,
	}
}

func TestNewMessageSentEvent(t *testing.T) {
	msg := createTestMessage()
	messageID := "msg_123"
	event := NewMessageSentEvent(msg, messageID)

	if event.Name != EventMessageSent {
		t.Errorf("Expected event name to be %s, got %s", EventMessageSent, event.Name)
	}

	if event.AggregateID != "123" {
		t.Errorf("Expected aggregate ID to be 123, got %s", event.AggregateID)
	}

	if event.Message != msg {
		t.Error("Expected message to be the same instance")
	}

	if event.MessageID != messageID {
		t.Errorf("Expected message ID to be %s, got %s", messageID, event.MessageID)
	}
}

func TestNewMessageFailedEvent(t *testing.T) {
	msg := createTestMessage()
	testError := errors.New("test error")
	event := NewMessageFailedEvent(msg, testError)

	if event.Name != EventMessageFailed {
		t.Errorf("Expected event name to be %s, got %s", EventMessageFailed, event.Name)
	}

	if event.AggregateID != "123" {
		t.Errorf("Expected aggregate ID to be 123, got %s", event.AggregateID)
	}

	if event.Message != msg {
		t.Error("Expected message to be the same instance")
	}

	if event.Error != "test error" {
		t.Errorf("Expected error message to be 'test error', got %s", event.Error)
	}
}

func TestNewMessageQueuedEvent(t *testing.T) {
	msg := createTestMessage()
	event := NewMessageQueuedEvent(msg)

	if event.Name != EventMessageQueued {
		t.Errorf("Expected event name to be %s, got %s", EventMessageQueued, event.Name)
	}

	if event.AggregateID != "123" {
		t.Errorf("Expected aggregate ID to be 123, got %s", event.AggregateID)
	}

	if event.Message != msg {
		t.Error("Expected message to be the same instance")
	}
}
