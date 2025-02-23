package scheduler

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ercancavusoglu/messaging/internal/domain/message"
)

type Scheduler struct {
	interval       time.Duration
	messageService *message.Service
	running        atomic.Bool
	stopChan       chan struct{}
	mu             sync.Mutex
}

func NewScheduler(messageService *message.Service, interval time.Duration) *Scheduler {
	return &Scheduler{
		interval:       interval,
		messageService: messageService,
		stopChan:       make(chan struct{}),
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	if !s.running.CompareAndSwap(false, true) {
		s.mu.Unlock()
		return fmt.Errorf("scheduler is already running")
	}
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	fmt.Println("[Scheduler] Starting...")
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	fmt.Println("[Scheduler] Ticker started")
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[Scheduler] Stopping due to context cancellation")
			s.running.Store(false)
			return ctx.Err()
		case <-s.stopChan:
			fmt.Println("[Scheduler] Stop signal received")
			s.running.Store(false)
			return nil
		case <-ticker.C:
			if !s.running.Load() {
				continue
			}

			messages, err := s.messageService.GetPendingMessages()
			if err != nil {
				fmt.Printf("[Scheduler] Error getting pending messages: %v\n", err)
				continue
			}

			fmt.Printf("[Scheduler] Found %d pending messages\n", len(messages))
			for _, msg := range messages {
				if msg.Status != message.StatusPending {
					continue
				}
				fmt.Printf("[Scheduler] Queueing message ID: %d, Content: %s\n", msg.ID, msg.Content)
				if err := s.messageService.QueueMessage(msg); err != nil {
					fmt.Printf("[Scheduler] Error queueing message: %v\n", err)
				}
			}
		}
	}
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running.Load() {
		fmt.Println("[Scheduler] Stopping...")
		close(s.stopChan)
		s.running.Store(false)
	}
}

func (s *Scheduler) IsRunning() bool {
	return s.running.Load()
}
