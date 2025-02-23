package mocks

import (
	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockMessageService struct {
	mock.Mock
}

func (m *MockMessageService) GetPendingMessages() ([]*domain.Message, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockMessageService) GetSendedMessages() ([]*domain.Message, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockMessageService) QueueMessage(msg *domain.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockMessageService) List() ([]*domain.Message, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockMessageService) Publish(msg *domain.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}
