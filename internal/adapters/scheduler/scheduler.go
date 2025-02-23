package scheduler

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ercancavusoglu/messaging/internal/domain"
	"github.com/ercancavusoglu/messaging/internal/ports"
)

type SchedulerStatus struct {
	IsRunning bool
	LastRun   time.Time
}

type SchedulerService struct {
	messageService ports.MessageService
	interval       time.Duration
	running        atomic.Bool
	stopChan       chan struct{}
	mu             sync.Mutex
	logger         ports.Logger
}

func NewSchedulerService(messageService ports.MessageService, interval time.Duration, logger ports.Logger) *SchedulerService {
	return &SchedulerService{
		messageService: messageService,
		interval:       interval,
		stopChan:       make(chan struct{}),
		logger:         logger,
	}
}

func (s *SchedulerService) Start(ctx context.Context) error {
	s.mu.Lock()
	if !s.running.CompareAndSwap(false, true) {
		s.mu.Unlock()
		return fmt.Errorf("scheduler is already running")
	}
	s.stopChan = make(chan struct{})
	s.mu.Unlock()

	s.logger.Info("[Scheduler] Starting...")
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.logger.Info("[Scheduler] Ticker started")
	for {
		select {
		case <-ctx.Done():
			s.logger.Info("[Scheduler] Stopping due to context cancellation")
			s.running.Store(false)
			return ctx.Err()
		case <-s.stopChan:
			s.logger.Info("[Scheduler] Stop signal received")
			s.running.Store(false)
			return nil
		case <-ticker.C:
			if !s.running.Load() {
				continue
			}

			messages, err := s.messageService.GetPendingMessages()
			if err != nil {
				s.logger.Errorf("[Scheduler] Error getting pending messages: %v", err)
				continue
			}

			s.logger.Infof("[Scheduler] Found %d pending messages", len(messages))
			for _, msg := range messages {
				if msg.Status != domain.StatusPending {
					continue
				}
				s.logger.Infof("[Scheduler] Queueing message ID: %d, Content: %s", msg.ID, msg.Content)
				if err := s.messageService.QueueMessage(msg); err != nil {
					s.logger.Errorf("[Scheduler] Error queueing message: %v", err)
				}
			}
		}
	}
}

func (s *SchedulerService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running.Load() {
		s.logger.Info("[Scheduler] Stopping...")
		close(s.stopChan)
		s.running.Store(false)
	}
}

func (s *SchedulerService) IsRunning() bool {
	return s.running.Load()
}
