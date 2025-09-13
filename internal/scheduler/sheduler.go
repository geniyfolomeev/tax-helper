package scheduler

import (
	"context"
	"tax-helper/internal/service/task"
	"time"
)

type Scheduler struct {
	service *task.TasksService
	ticker  *time.Ticker
}

func NewScheduler(service *task.TasksService, interval time.Duration) *Scheduler {
	return &Scheduler{
		service: service,
		ticker:  time.NewTicker(interval),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.ticker.Stop()
				return
			default:
				err := s.service.ProcessDueReminders(ctx)
				if err != nil {

				}
			}
		}
	}()
}
