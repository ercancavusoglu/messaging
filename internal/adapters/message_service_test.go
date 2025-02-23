package adapters

import (
	"testing"
	"time"

	"github.com/ercancavusoglu/messaging/internal/adapters/mocks"
	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func createTestMessage() *domain.Message {
	return &domain.Message{
		ID:        123,
		To:        "+905551234567",
		Content:   "Test message",
		Status:    domain.StatusPending,
		MessageID: "",
		Provider:  "webhook",
		CreatedAt: time.Now(),
		SentAt:    nil,
	}
}

func TestMessageService_GetPendingMessages(t *testing.T) {
	mockRepo := &mocks.MockRepository{}
	mockWebhook := &mocks.MockWebhookClient{}
	mockCache := &mocks.MockCache{}
	mockEventBus := &mocks.MockEventBus{}

	service := NewMessageService(mockRepo, mockWebhook, mockCache, mockEventBus)

	expectedMessages := []*domain.Message{createTestMessage()}
	mockRepo.On("GetByStatus", domain.StatusPending).Return(expectedMessages, nil)

	messages, err := service.GetPendingMessages()
	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}

func TestMessageService_GetSendedMessages(t *testing.T) {
	mockRepo := &mocks.MockRepository{}
	mockWebhook := &mocks.MockWebhookClient{}
	mockCache := &mocks.MockCache{}
	mockEventBus := &mocks.MockEventBus{}

	service := NewMessageService(mockRepo, mockWebhook, mockCache, mockEventBus)

	expectedMessages := []*domain.Message{createTestMessage()}
	mockRepo.On("GetByStatus", domain.StatusSent).Return(expectedMessages, nil)

	messages, err := service.GetSendedMessages()
	assert.NoError(t, err)
	assert.Equal(t, expectedMessages, messages)
	mockRepo.AssertExpectations(t)
}

func TestMessageService_Publish(t *testing.T) {
	mockRepo := &mocks.MockRepository{}
	mockWebhook := &mocks.MockWebhookClient{}
	mockCache := &mocks.MockCache{}
	mockEventBus := &mocks.MockEventBus{}

	service := NewMessageService(mockRepo, mockWebhook, mockCache, mockEventBus)

	msg := createTestMessage()
	mockEventBus.On("Publish", mock.AnythingOfType("domain.MessageQueuedEvent")).Return(nil)

	err := service.Publish(msg)
	assert.NoError(t, err)
	mockEventBus.AssertExpectations(t)
}
