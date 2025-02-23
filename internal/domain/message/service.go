package message

import (
	"fmt"
	"time"

	"github.com/ercancavusoglu/messaging/internal/domain/event"
)

type Service struct {
	repo          Repository
	webhookClient WebhookClient
	cache         Cache
	eventBus      event.EventBus
}

func NewService(repo Repository, webhookClient WebhookClient, cache Cache, eventBus event.EventBus) *Service {
	return &Service{
		repo:          repo,
		webhookClient: webhookClient,
		cache:         cache,
		eventBus:      eventBus,
	}
}

func (s *Service) Create(to, content string) (*Message, error) {
	msg := &Message{
		To:        to,
		Content:   content,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(msg); err != nil {
		return nil, fmt.Errorf("failed to create message: %v", err)
	}

	s.eventBus.Publish(NewMessageCreatedEvent(msg))

	return msg, nil
}

func (s *Service) GetByID(id int64) (*Message, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetPendingMessages() ([]*Message, error) {
	return s.repo.GetByStatus(StatusPending)
}

func (s *Service) QueueMessage(msg *Message) error {
	s.eventBus.Publish(NewMessageQueuedEvent(msg))
	return nil
}

func (s *Service) List(limit int) ([]*Message, error) {
	return s.repo.GetPendingMessages(limit)
}

func (s *Service) SendMessage(msg *Message) error {
	resp, err := s.webhookClient.SendMessage(msg.To, msg.Content)
	if err != nil {
		s.repo.UpdateStatus(msg.ID, StatusFailed, "")
		s.eventBus.Publish(NewMessageFailedEvent(msg, err))
		return err
	}

	now := time.Now()
	s.repo.UpdateStatus(msg.ID, StatusSent, resp.MessageID)
	s.cache.Set(resp.MessageID, now)

	s.eventBus.Publish(NewMessageSentEvent(msg, resp.MessageID))
	return nil
}
