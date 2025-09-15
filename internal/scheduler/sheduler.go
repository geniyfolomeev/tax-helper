package scheduler

import (
	"context"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"tax-helper/internal/infrastructure/db"
	"tax-helper/internal/infrastructure/processors"
	"tax-helper/internal/infrastructure/repository"
	"tax-helper/internal/logger"
	"tax-helper/internal/service/task"
)

type Notifier interface {
	Process(ctx context.Context, task db.Tasks) error
}

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (c *RealClock) Now() time.Time {
	return time.Now()
}

type TxManager interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type Scheduler struct {
	svc        task.TasksService
	ticker     *time.Ticker
	processors map[string]Notifier
	logger     logger.Logger
	repo       repository.TasksRepository
	txManager  TxManager
	clock      Clock
}

func NewScheduler(
	taskService task.TasksService,
	interval time.Duration,
	botClient *tgbotapi.BotAPI,
	logger logger.Logger,
	repo repository.TasksRepository,
	txManager TxManager,
	clock Clock,
) *Scheduler {
	s := &Scheduler{
		svc:       taskService,
		ticker:    time.NewTicker(interval),
		logger:    logger,
		repo:      repo,
		txManager: txManager,
		clock:     clock,
	}

	s.processors = map[string]Notifier{
		"submit_declaration": processors.NewSendDeclarationProcessor(botClient, taskService, logger, repo, txManager),
		"add_income":         processors.NewAddIncomeProcessor(botClient, taskService, logger, repo, txManager),
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
			s.logger.Info("tick: fetching  tasks")

			now := s.clock.Now()
			tasks, err := s.svc.GetDueTasks(ctx, 100, now)
			if err != nil {
				s.logger.Errorf("error getting tasks: %v", err)
				continue
			}
			for _, t := range tasks {
				if proc, ok := s.processors[t.Type]; ok {
					if err := proc.Process(ctx, t); err != nil {
						s.logger.Errorf("processor error for task %v: %v", t, err)
					}
				} else {
					s.logger.Error("no processor for task type %q", t.Type)
				}
			}
		}
	}
}
