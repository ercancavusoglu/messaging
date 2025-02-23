package adapters

import (
	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/ercancavusoglu/messaging/internal/ports"
)

type messageService struct {
	repo          ports.Repository
	webhookClient ports.WebhookClient
	cache         ports.Cache
	eventBus      ports.EventBus
}

func NewMessageService(repo ports.Repository, webhookClient ports.WebhookClient, cache ports.Cache, eventBus ports.EventBus) ports.MessageService {
	return &messageService{
		repo:          repo,
		webhookClient: webhookClient,
		cache:         cache,
		eventBus:      eventBus,
	}
}

func (s *messageService) GetPendingMessages() ([]*domain.Message, error) {
	return s.repo.GetByStatus(domain.StatusPending)
}

func (s *messageService) GetSendedMessages() ([]*domain.Message, error) {
	return s.repo.GetByStatus(domain.StatusSent)
}

func (s *messageService) QueueMessage(msg *domain.Message) error {
	s.eventBus.Publish(domain.NewMessageQueuedEvent(msg))
	return nil
}

func (s *messageService) List() ([]*domain.Message, error) {
	return s.repo.GetByStatus(domain.StatusSent)
}
