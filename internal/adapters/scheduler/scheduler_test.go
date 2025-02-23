package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/ercancavusoglu/messaging/internal/adapters/mocks"
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

func createTestMessage() *domain.Message {
	return &domain.Message{
		ID:        123,
		To:        "+905551234567",
		Content:   "Test message",
		Status:    domain.StatusPending,
		MessageID: "",
		CreatedAt: time.Now(),
		SentAt:    nil,
	}
}

func TestSchedulerService_Start_Success(t *testing.T) {
	mockService := &mocks.MockMessageService{}
	logger := &mockLogger{}

	scheduler := NewSchedulerService(mockService, 100*time.Millisecond, logger)

	// Mock beklentileri
	msg := createTestMessage()
	messages := []*domain.Message{msg}
	mockService.On("GetPendingMessages").Return(messages, nil)
	mockService.On("Publish", msg).Return(nil)

	// Scheduler'ı başlat
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = scheduler.Start(ctx)
	}()

	// Scheduler'ın çalışması için bekle
	time.Sleep(200 * time.Millisecond)

	// Scheduler'ı durdur
	cancel()

	// Beklentilerin karşılandığını kontrol et
	mockService.AssertExpectations(t)
}

func TestSchedulerService_Start_AlreadyRunning(t *testing.T) {
	mockService := &mocks.MockMessageService{}
	logger := &mockLogger{}

	scheduler := NewSchedulerService(mockService, 100*time.Millisecond, logger)

	// İlk başlatma
	ctx1, cancel1 := context.WithCancel(context.Background())
	go func() {
		_ = scheduler.Start(ctx1)
	}()
	time.Sleep(50 * time.Millisecond) // Scheduler'ın başlaması için bekle

	// İkinci başlatma denemesi
	ctx2 := context.Background()
	err := scheduler.Start(ctx2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "scheduler is already running")

	// Temizlik
	cancel1()
	scheduler.Stop()
}

func TestSchedulerService_Stop(t *testing.T) {
	mockService := &mocks.MockMessageService{}
	logger := &mockLogger{}

	scheduler := NewSchedulerService(mockService, 100*time.Millisecond, logger)

	// Scheduler'ı başlat
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = scheduler.Start(ctx)
	}()
	time.Sleep(50 * time.Millisecond) // Scheduler'ın başlaması için bekle

	// Scheduler'ı durdur
	scheduler.Stop()
	assert.False(t, scheduler.IsRunning())
}

func TestSchedulerService_IsRunning(t *testing.T) {
	mockService := &mocks.MockMessageService{}
	logger := &mockLogger{}

	scheduler := NewSchedulerService(mockService, 100*time.Millisecond, logger)

	// Başlangıçta çalışmıyor olmalı
	assert.False(t, scheduler.IsRunning())

	// Scheduler'ı başlat
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = scheduler.Start(ctx)
	}()
	time.Sleep(50 * time.Millisecond) // Scheduler'ın başlaması için bekle

	// Çalışıyor olmalı
	assert.True(t, scheduler.IsRunning())

	// Scheduler'ı durdur
	scheduler.Stop()

	// Tekrar çalışmıyor olmalı
	assert.False(t, scheduler.IsRunning())
}
