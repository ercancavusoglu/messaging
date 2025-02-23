package scheduler

import "time"

type SchedulerStatus struct {
	IsRunning bool
	LastRun   time.Time
}
