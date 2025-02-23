package scheduler

import (
	"context"
	"sync"
	"time"
)

type Task interface {
	Execute() error
}

type Scheduler struct {
	interval time.Duration
	tasks    []Task
	running  bool
	mu       sync.Mutex
	stop     chan struct{}
}

func NewScheduler(interval time.Duration) *Scheduler {
	return &Scheduler{
		interval: interval,
		tasks:    make([]Task, 0),
		stop:     make(chan struct{}),
	}
}

func (s *Scheduler) AddTask(task Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, task)
}

func (s *Scheduler) Start(ctx context.Context) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = true
	s.mu.Unlock()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-s.stop:
			return nil
		case <-ticker.C:
			s.executeTasks()
		}
	}
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		s.running = false
		close(s.stop)
	}
}

func (s *Scheduler) executeTasks() {
	s.mu.Lock()
	tasks := s.tasks
	s.mu.Unlock()

	for _, task := range tasks {
		if err := task.Execute(); err != nil {
			// Log error but continue with other tasks
			continue
		}
	}
}
