package adapters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ercancavusoglu/messaging/internal/adapters/mocks"
	"github.com/ercancavusoglu/messaging/internal/adapters/scheduler"
	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/ercancavusoglu/messaging/internal/ports"
	"github.com/stretchr/testify/assert"
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

func TestMessageHandler_GetMessages(t *testing.T) {
	mockService := &mocks.MockMessageService{}
	mockScheduler := scheduler.NewSchedulerService(mockService, time.Second, &mockLogger{})

	handler := NewMessageHandler(mockService, mockScheduler)

	expectedMessages := []*domain.Message{createTestMessage()}
	mockService.On("GetSendedMessages").Return(expectedMessages, nil)

	req := httptest.NewRequest(http.MethodGet, "/messages", nil)
	w := httptest.NewRecorder()

	handler.GetMessages(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []*domain.Message
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)

	// CreatedAt alanını karşılaştırmadan önce sıfırlayalım
	for _, msg := range response {
		msg.CreatedAt = expectedMessages[0].CreatedAt
	}
	assert.Equal(t, expectedMessages, response)

	mockService.AssertExpectations(t)
}

func TestMessageHandler_StartScheduler_AlreadyRunning(t *testing.T) {
	mockService := &mocks.MockMessageService{}
	mockScheduler := scheduler.NewSchedulerService(mockService, time.Second, &mockLogger{})

	handler := NewMessageHandler(mockService, mockScheduler)

	// Scheduler'ı başlatalım
	go func() {
		_ = mockScheduler.Start(handler.ctx)
	}()
	time.Sleep(time.Millisecond * 100) // Scheduler'ın başlaması için bekleyelim

	req := httptest.NewRequest(http.MethodPost, "/scheduler/start", nil)
	w := httptest.NewRecorder()

	handler.StartScheduler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Scheduler is already running", response["error"])

	mockScheduler.Stop()
}

func TestMessageHandler_StopScheduler_NotRunning(t *testing.T) {
	mockService := &mocks.MockMessageService{}
	mockScheduler := scheduler.NewSchedulerService(mockService, time.Second, &mockLogger{})

	handler := NewMessageHandler(mockService, mockScheduler)

	req := httptest.NewRequest(http.MethodPost, "/scheduler/stop", nil)
	w := httptest.NewRecorder()

	handler.StopScheduler(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Scheduler is not running", response["error"])
}
