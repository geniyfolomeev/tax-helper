package scheduler

import (
	"context"
	"tax-helper/internal/logger"
	"tax-helper/internal/service"
	"time"
)

type Scheduler struct {
	service service.TaskService
	ticker  *time.Ticker
}

func NewScheduler(service service.TaskService, interval time.Duration) *Scheduler {
	return &Scheduler{
		service: service,
		ticker:  time.NewTicker(interval),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	go func() {
		logger.Info("Scheduler started")
		for {
			select {
			case <-ctx.Done():
				s.ticker.Stop()
				logger.Info("Scheduler stopped")
				return
			default:
				logger.Info("Scheduler try to send messages")
				err := s.service.ProcessDueReminders(ctx)
				if err != nil {
					logger.Error("Scheduler process due to error:", err)

				}
			}
		}
	}()
}
