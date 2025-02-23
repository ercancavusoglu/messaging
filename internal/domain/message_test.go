package domain

import (
	"testing"
	"time"
)

func TestMessageStatus_Constants(t *testing.T) {
	tests := []struct {
		name     string
		status   MessageStatus
		expected string
	}{
		{
			name:     "Pending status",
			status:   StatusPending,
			expected: "pending",
		},
		{
			name:     "Sent status",
			status:   StatusSent,
			expected: "sent",
		},
		{
			name:     "Failed status",
			status:   StatusFailed,
			expected: "failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.expected {
				t.Errorf("MessageStatus %s = %v, want %v", tt.name, tt.status, tt.expected)
			}
		})
	}
}

func TestMessage_Creation(t *testing.T) {
	now := time.Now()
	msg := Message{
		ID:        1,
		To:        "+905551234567",
		Content:   "Test message",
		Status:    StatusPending,
		MessageID: "msg_123",
		Provider:  "webhook",
		CreatedAt: now,
		SentAt:    nil,
	}

	if msg.ID != 1 {
		t.Errorf("Expected ID to be 1, got %d", msg.ID)
	}

	if msg.To != "+905551234567" {
		t.Errorf("Expected To to be +905551234567, got %s", msg.To)
	}

	if msg.Content != "Test message" {
		t.Errorf("Expected Content to be 'Test message', got %s", msg.Content)
	}

	if msg.Status != StatusPending {
		t.Errorf("Expected Status to be pending, got %s", msg.Status)
	}

	if msg.MessageID != "msg_123" {
		t.Errorf("Expected MessageID to be msg_123, got %s", msg.MessageID)
	}

	if msg.Provider != "webhook" {
		t.Errorf("Expected Provider to be webhook, got %s", msg.Provider)
	}

	if !msg.CreatedAt.Equal(now) {
		t.Errorf("Expected CreatedAt to be %v, got %v", now, msg.CreatedAt)
	}

	if msg.SentAt != nil {
		t.Error("Expected SentAt to be nil")
	}
}
