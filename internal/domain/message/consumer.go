package message

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/ercancavusoglu/messaging/internal/domain/event"
)

type Consumer struct {
	webhookClient WebhookClient
	repo          Repository
	cache         Cache
	eventBus      event.EventBus
	workers       int
	workerPool    chan struct{}
	wg            sync.WaitGroup
}

func NewConsumer(webhookClient WebhookClient, repo Repository, cache Cache, eventBus event.EventBus, workers int) *Consumer {
	return &Consumer{
		webhookClient: webhookClient,
		repo:          repo,
		cache:         cache,
		eventBus:      eventBus,
		workers:       workers,
		workerPool:    make(chan struct{}, workers),
	}
}

func (c *Consumer) Start() error {
	fmt.Println("[Consumer] Starting...")

	c.eventBus.Subscribe(EventMessageQueued, func(e event.Event) error {
		var evt MessageQueuedEvent
		if err := json.Unmarshal(e.(*event.EventEnvelope).Data, &evt); err != nil {
			fmt.Printf("[Consumer] Failed to unmarshal event: %v\n", err)
			return fmt.Errorf("failed to unmarshal event: %v", err)
		}

		msg := evt.Message
		fmt.Printf("[Consumer] Processing queued message ID: %d, To: %s\n", msg.ID, msg.To)

		c.workerPool <- struct{}{}
		c.wg.Add(1)

		go func() {
			defer func() {
				<-c.workerPool
				c.wg.Done()
			}()

			if err := c.processMessage(msg); err != nil {
				fmt.Printf("[Consumer] Failed to process message ID: %d, Error: %v\n", msg.ID, err)
			}
		}()

		return nil
	})

	return nil
}

func (c *Consumer) Stop() {
	c.wg.Wait()
}

func (c *Consumer) processMessage(msg *Message) error {
	log.Printf("[Consumer] Processing message [id: %d]", msg.ID)

	webhookResponse, err := c.webhookClient.SendMessage(msg.To, msg.Content)
	if err != nil {
		log.Printf("[Consumer] Failed to send message to webhook: %v", err)
		if err := c.repo.UpdateStatus(msg.ID, StatusFailed, ""); err != nil {
			log.Printf("[Consumer] Failed to update message status: %v", err)
		}
		return fmt.Errorf("failed to send message to webhook: %v", err)
	}

	if err := c.repo.UpdateStatus(msg.ID, StatusSent, webhookResponse.MessageID); err != nil {
		log.Printf("[Consumer] Failed to update message status: %v", err)
		return fmt.Errorf("failed to update message status: %v", err)
	}

	cacheKey := fmt.Sprintf("message:%d", msg.ID)
	if err := c.cache.Set(cacheKey, nil); err != nil {
		log.Printf("[Consumer] Failed to delete message from cache: %v", err)
	}

	event := NewMessageSentEvent(msg, webhookResponse.MessageID)

	if err := c.eventBus.Publish(&event); err != nil {
		log.Printf("[Consumer] Failed to publish message sent event: %v", err)
	}

	log.Printf("[Consumer] Message processed successfully [id: %d]", msg.ID)
	return nil
}
