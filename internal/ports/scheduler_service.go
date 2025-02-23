package ports

import (
	"context"
)

type SchedulerService interface {
	Start(ctx context.Context) error
	Stop()
	IsRunning() bool
}
