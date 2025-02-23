package message

import (
	"time"
)

type WebhookClient interface {
	SendMessage(to, content string) (*WebhookResponse, error)
}

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
}

type Repository interface {
	List() ([]*Message, error)
	GetPendingMessages(limit int) ([]*Message, error)
	UpdateStatus(id int64, status MessageStatus, messageID string) error
}

type Service struct {
	repo          Repository
	webhookClient WebhookClient
	cache         Cache
}

func NewService(repo Repository, webhookClient WebhookClient, cache Cache) *Service {
	return &Service{
		repo:          repo,
		webhookClient: webhookClient,
		cache:         cache,
	}
}

func (s *Service) SendPendingMessages() error {
	messages, err := s.repo.GetPendingMessages(2)
	if err != nil {
		return err
	}

	for _, msg := range messages {
		if len(msg.Content) > MaxContentLen {

			continue
		}

		messageID, err := s.webhookClient.SendMessage(msg.To, msg.Content)
		if err != nil {
			s.repo.UpdateStatus(msg.ID, StatusFailed, "")
			continue
		}

		now := time.Now()
		s.repo.UpdateStatus(msg.ID, StatusSent, messageID.MessageID)

		s.cache.Set(messageID.MessageID, now)
	}

	return nil
}

func (s *Service) List() ([]*Message, error) {
	return s.repo.List()
}

func (s *Service) ProcessMessage(msg *Message) error {
	response, err := s.webhookClient.SendMessage(msg.To, msg.Content)
	if err != nil {
		err = s.repo.UpdateStatus(msg.ID, StatusFailed, "")
		return err
	}

	err = s.repo.UpdateStatus(msg.ID, StatusSent, response.MessageID)
	if err != nil {
		return err
	}

	// Cache the message ID and sent time
	err = s.cache.Set(response.MessageID, time.Now().String())
	return err
}

type WebhookResponse struct {
	Message   string `json:"message"`
	MessageID string `json:"messageId"`
}
