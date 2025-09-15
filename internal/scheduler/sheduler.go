// internal/scheduler/scheduler.go
package scheduler

import (
	"context"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/bot"
	"tax-helper/internal/infrastructure/processors"
	"tax-helper/internal/infrastructure/repository"
	"tax-helper/internal/logger"
	"tax-helper/internal/service/task"
	"time"
)

type Scheduler struct {
	svc        task.TasksService
	ticker     *time.Ticker
	processors map[string]domain.Notifier
	logger     logger.Logger
	repo       repository.TasksRepo
}

func NewScheduler(svc task.TasksService, interval time.Duration, botClient *bot.Bot, logger logger.Logger, repo repository.TasksRepo) *Scheduler {
	s := &Scheduler{
		svc:    svc,
		ticker: time.NewTicker(interval),
		logger: logger,
		repo:   repo,
	}

	// Составляем мапу процессоров сразу здесь
	s.processors = map[string]domain.Notifier{
		"send_declaration": processors.NewSendDeclarationProcessor(botClient, svc, logger, repo),
		"add_income":       processors.NewAddIncomeProcessor(botClient, svc, logger, repo),
	}

	return s
}

func (s *Scheduler) Start(ctx context.Context) {
	s.logger.Info("scheduler started")
	defer s.ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Info("scheduler stopped by context")
			return
		case <-s.ticker.C:
			s.logger.Info("tick: fetching pending tasks")
			tasks, err := s.svc.GetDueTasks(ctx)
			if err != nil {
				s.logger.Info("error getting tasks: %v\n", err)
				continue
			}

			for _, t := range tasks {
				if proc, ok := s.processors[t.Type]; ok {
					if err := proc.Process(ctx, t); err != nil {
						s.logger.Info("processor error for task %v: %v\n", t, err)
					}
				} else {
					s.logger.Info("no processor for task type %q\n", t.Type)
				}
			}
		}
	}
}
