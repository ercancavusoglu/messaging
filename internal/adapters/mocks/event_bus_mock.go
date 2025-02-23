package mocks

import (
	"github.com/ercancavusoglu/messaging/internal/ports"
	"github.com/stretchr/testify/mock"
)

type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Publish(event ports.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

func (m *MockEventBus) Subscribe(eventName string, handler ports.EventHandler) {
	m.Called(eventName, handler)
}

func (m *MockEventBus) Unsubscribe(eventName string, handler ports.EventHandler) {
	m.Called(eventName, handler)
}
