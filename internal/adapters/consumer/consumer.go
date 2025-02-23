package consumer

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/ercancavusoglu/messaging/internal/ports"
)

type Consumer struct {
	webhookClient ports.WebhookClient
	repo          ports.Repository
	cache         ports.Cache
	eventBus      ports.EventBus
	workers       int
	workerPool    chan struct{}
	wg            sync.WaitGroup
	logger        ports.Logger
}

func NewConsumer(webhookClient ports.WebhookClient, repo ports.Repository, cache ports.Cache, eventBus ports.EventBus, workers int, logger ports.Logger) *Consumer {
	return &Consumer{
		webhookClient: webhookClient,
		repo:          repo,
		cache:         cache,
		eventBus:      eventBus,
		workers:       workers,
		workerPool:    make(chan struct{}, workers),
		logger:        logger,
	}
}

func (c *Consumer) Start() error {
	c.logger.Info("[Consumer] Starting...")

	c.eventBus.Subscribe(domain.EventMessageQueued, func(e ports.Event) error {
		var evt domain.MessageQueuedEvent
		if err := json.Unmarshal(e.(*domain.EventEnvelope).Data, &evt); err != nil {
			c.logger.Errorf("[Consumer] Failed to unmarshal event: %v", err)
			return fmt.Errorf("failed to unmarshal event: %v", err)
		}

		msg := evt.Message
		c.logger.Infof("[Consumer] Processing queued message ID: %d, To: %s", msg.ID, msg.To)

		c.workerPool <- struct{}{}
		c.wg.Add(1)

		go func() {
			defer func() {
				<-c.workerPool
				c.wg.Done()
			}()

			if err := c.processMessage(msg); err != nil {
				c.logger.Errorf("[Consumer] Failed to process message ID: %d, Error: %v", msg.ID, err)
			}
		}()

		return nil
	})

	return nil
}

func (c *Consumer) Stop() {
	c.wg.Wait()
}

func (c *Consumer) processMessage(msg *domain.Message) error {
	c.logger.Infof("[Consumer] Processing message [id: %d]", msg.ID)

	webhookResponse, err := c.webhookClient.SendMessage(msg.To, msg.Content)
	if err != nil {
		c.logger.Errorf("[Consumer] Failed to send message to webhook: %v", err)
		if err := c.repo.UpdateStatus(msg.ID, domain.StatusFailed, "", ""); err != nil {
			c.logger.Errorf("[Consumer] Failed to update message status: %v", err)
		}
		return fmt.Errorf("failed to send message to webhook: %v", err)
	}

	if err := c.repo.UpdateStatus(msg.ID, domain.StatusSent, webhookResponse.MessageID, webhookResponse.Provider); err != nil {
		c.logger.Errorf("[Consumer] Failed to update message status: %v", err)
		return fmt.Errorf("failed to update message status: %v", err)
	}

	cacheKey := fmt.Sprintf("message:%d", msg.ID)
	if err := c.cache.Set(cacheKey, nil); err != nil {
		c.logger.Errorf("[Consumer] Failed to delete message from cache: %v", err)
	}

	event := domain.NewMessageSentEvent(msg, webhookResponse.MessageID)

	if err := c.eventBus.Publish(&event); err != nil {
		c.logger.Errorf("[Consumer] Failed to publish message sent event: %v", err)
	}

	c.logger.Infof("[Consumer] Message processed successfully [id: %d]", msg.ID)
	return nil
}
