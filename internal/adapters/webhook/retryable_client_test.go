package webhook

import (
	"errors"
	"testing"

	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/ercancavusoglu/messaging/internal/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWebhookClient struct {
	mock.Mock
}

func (m *MockWebhookClient) SendMessage(to, content string) (*domain.WebhookResponse, error) {
	args := m.Called(to, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.WebhookResponse), args.Error(1)
}

func TestRetryableWebhookClient_SendMessage_Success(t *testing.T) {
	mockClient1 := new(MockWebhookClient)
	mockClient2 := new(MockWebhookClient)

	expectedResponse := &domain.WebhookResponse{
		MessageID: "msg_123",
		Message:   "Message sent successfully",
	}

	// İlk client hata döndürür
	mockClient1.On("SendMessage", "+905551234567", "Test message").Return(nil, errors.New("connection error"))

	// İkinci client başarılı olur
	mockClient2.On("SendMessage", "+905551234567", "Test message").Return(expectedResponse, nil)

	clients := []ports.WebhookClient{mockClient1, mockClient2}
	retryableClient := NewRetryableWebhookClient(clients, 3)

	response, err := retryableClient.SendMessage("+905551234567", "Test message")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedResponse.MessageID, response.MessageID)
	assert.Equal(t, expectedResponse.Message, response.Message)

	mockClient1.AssertExpectations(t)
	mockClient2.AssertExpectations(t)
}

func TestRetryableWebhookClient_SendMessage_AllFail(t *testing.T) {
	mockClient1 := new(MockWebhookClient)
	mockClient2 := new(MockWebhookClient)

	// Her iki client de hata döndürür
	mockClient1.On("SendMessage", "+905551234567", "Test message").Return(nil, errors.New("connection error"))
	mockClient2.On("SendMessage", "+905551234567", "Test message").Return(nil, errors.New("timeout error"))

	clients := []ports.WebhookClient{mockClient1, mockClient2}
	retryableClient := NewRetryableWebhookClient(clients, 3)

	response, err := retryableClient.SendMessage("+905551234567", "Test message")
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "all retry attempts failed")

	mockClient1.AssertExpectations(t)
	mockClient2.AssertExpectations(t)
}

func TestRetryableWebhookClient_SendMessage_EventualSuccess(t *testing.T) {
	mockClient := new(MockWebhookClient)

	expectedResponse := &domain.WebhookResponse{
		MessageID: "msg_123",
		Message:   "Message sent successfully",
	}

	// İlk iki deneme başarısız, üçüncü deneme başarılı
	mockClient.On("SendMessage", "+905551234567", "Test message").
		Return(nil, errors.New("error 1")).Once()
	mockClient.On("SendMessage", "+905551234567", "Test message").
		Return(nil, errors.New("error 2")).Once()
	mockClient.On("SendMessage", "+905551234567", "Test message").
		Return(expectedResponse, nil).Once()

	clients := []ports.WebhookClient{mockClient}
	retryableClient := NewRetryableWebhookClient(clients, 3)

	response, err := retryableClient.SendMessage("+905551234567", "Test message")
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, expectedResponse.MessageID, response.MessageID)
	assert.Equal(t, expectedResponse.Message, response.Message)

	mockClient.AssertExpectations(t)
}
