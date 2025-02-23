package scheduler

import (
	"context"
	"fmt"
	"time"
)

type MessageService interface {
	SendPendingMessages() error
}

type Scheduler struct {
	messageService MessageService
	interval       time.Duration
	running        bool
}

func NewScheduler(messageService MessageService, interval time.Duration) *Scheduler {
	return &Scheduler{
		messageService: messageService,
		interval:       interval,
		running:        false,
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.running = true
	fmt.Println("Scheduler started")
	go func() {
		for s.running {
			s.messageService.SendPendingMessages()
			time.Sleep(s.interval)
			fmt.Println("Scheduler running")
			select {
			case <-ctx.Done():
				fmt.Println("Scheduler stopped")
				s.running = false
				return
			default:
				fmt.Println("Scheduler running...")
				continue
			}

		}
	}()

	return nil
}

func (s *Scheduler) Stop() {
	s.running = false
	fmt.Println("Scheduler stopped")
}

func (s *Scheduler) IsRunning() bool {
	return s.running
}
