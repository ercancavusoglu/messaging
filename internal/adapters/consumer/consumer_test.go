package consumer

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/ercancavusoglu/messaging/internal/adapters/mocks"
	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/ercancavusoglu/messaging/internal/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLogger struct {
	ports.Logger
}

func (m *mockLogger) Info(args ...interface{})                    {}
func (m *mockLogger) Infof(format string, args ...interface{})    {}
func (m *mockLogger) Error(args ...interface{})                   {}
func (m *mockLogger) Errorf(format string, args ...interface{})   {}
func (m *mockLogger) Warning(args ...interface{})                 {}
func (m *mockLogger) Warningf(format string, args ...interface{}) {}

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

func TestConsumer_ProcessMessage_Success(t *testing.T) {
	mockWebhook := &mocks.MockWebhookClient{}
	mockRepo := &mocks.MockRepository{}
	mockCache := &mocks.MockCache{}
	mockEventBus := &mocks.MockEventBus{}
	logger := &mockLogger{}

	consumer := NewConsumer(mockWebhook, mockRepo, mockCache, mockEventBus, 1, logger)

	msg := createTestMessage()
	webhookResponse := &domain.WebhookResponse{
		MessageID: "msg_123",
		Message:   "Message sent successfully",
		Provider:  "client_one",
	}

	// Mock beklentileri
	mockWebhook.On("SendMessage", msg.To, msg.Content).Return(webhookResponse, nil)
	mockRepo.On("UpdateStatus", msg.ID, domain.StatusSent, webhookResponse.MessageID, webhookResponse.Provider).Return(nil)
	mockCache.On("Set", mock.Anything, mock.Anything).Return(nil)
	mockEventBus.On("Publish", mock.AnythingOfType("*domain.MessageSentEvent")).Return(nil)

	// Test
	err := consumer.processMessage(msg)
	assert.NoError(t, err)

	// Beklentilerin karşılandığını kontrol et
	mockWebhook.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

func TestConsumer_ProcessMessage_WebhookError(t *testing.T) {
	mockWebhook := &mocks.MockWebhookClient{}
	mockRepo := &mocks.MockRepository{}
	mockCache := &mocks.MockCache{}
	mockEventBus := &mocks.MockEventBus{}
	logger := &mockLogger{}

	consumer := NewConsumer(mockWebhook, mockRepo, mockCache, mockEventBus, 1, logger)

	msg := createTestMessage()

	// Mock beklentileri
	mockWebhook.On("SendMessage", msg.To, msg.Content).Return(nil, assert.AnError)
	mockRepo.On("UpdateStatus", msg.ID, domain.StatusFailed, "", "").Return(nil)

	// Test
	err := consumer.processMessage(msg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to send message to webhook")

	// Beklentilerin karşılandığını kontrol et
	mockWebhook.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestConsumer_Start(t *testing.T) {
	mockWebhook := &mocks.MockWebhookClient{}
	mockRepo := &mocks.MockRepository{}
	mockCache := &mocks.MockCache{}
	mockEventBus := &mocks.MockEventBus{}
	logger := &mockLogger{}

	consumer := NewConsumer(mockWebhook, mockRepo, mockCache, mockEventBus, 1, logger)

	// Mock event handler'ı yakalayalım
	var capturedHandler ports.EventHandler
	mockEventBus.On("Subscribe", domain.EventMessageQueued, mock.AnythingOfType("ports.EventHandler")).
		Run(func(args mock.Arguments) {
			capturedHandler = args.Get(1).(ports.EventHandler)
		}).
		Return()

	// Consumer'ı başlat
	err := consumer.Start()
	assert.NoError(t, err)

	// Test mesajı oluştur
	msg := createTestMessage()
	event := domain.NewMessageQueuedEvent(msg)

	// Event envelope oluştur
	data, err := json.Marshal(event)
	assert.NoError(t, err)

	envelope := &domain.EventEnvelope{
		Name:        event.EventName(),
		OccurredOn:  event.OccurredAt(),
		AggregateID: event.GetAggregateID(),
		Data:        data,
	}

	// Webhook response'u hazırla
	webhookResponse := &domain.WebhookResponse{
		MessageID: "msg_123",
		Message:   "Message sent successfully",
		Provider:  "client_one",
	}

	// Mock beklentileri
	mockWebhook.On("SendMessage", msg.To, msg.Content).Return(webhookResponse, nil)
	mockRepo.On("UpdateStatus", msg.ID, domain.StatusSent, webhookResponse.MessageID, webhookResponse.Provider).Return(nil)
	mockCache.On("Set", mock.Anything, mock.Anything).Return(nil)
	mockEventBus.On("Publish", mock.AnythingOfType("*domain.MessageSentEvent")).Return(nil)

	// Event handler'ı çağır
	err = capturedHandler(envelope)
	assert.NoError(t, err)

	// İşlemin tamamlanmasını bekle
	consumer.Stop()

	// Beklentilerin karşılandığını kontrol et
	mockWebhook.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	mockCache.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}
