package mocks

import (
	"github.com/ercancavusoglu/messaging/internal/domain"
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
