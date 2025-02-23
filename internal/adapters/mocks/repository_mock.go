package mocks

import (
	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(message *domain.Message) error {
	args := m.Called(message)
	return args.Error(0)
}

func (m *MockRepository) GetPendingMessages() ([]*domain.Message, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockRepository) UpdateStatus(id int64, status domain.MessageStatus, messageID string, provider string) error {
	args := m.Called(id, status, messageID, provider)
	return args.Error(0)
}

func (m *MockRepository) GetByStatus(status domain.MessageStatus) ([]*domain.Message, error) {
	args := m.Called(status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Message), args.Error(1)
}

func (m *MockRepository) UpdateMessageID(id int64, messageID string) error {
	args := m.Called(id, messageID)
	return args.Error(0)
}
