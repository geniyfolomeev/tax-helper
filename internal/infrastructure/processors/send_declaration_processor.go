package processors

import (
	"context"
	"fmt"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/bot"
	"tax-helper/internal/infrastructure/db"
	"tax-helper/internal/infrastructure/repository"
	"tax-helper/internal/logger"
	"tax-helper/internal/service/task"
	"time"
)

type SendDeclarationProcessor struct {
	bot    *bot.Bot
	tasksS task.TasksService
	logger logger.Logger
	repo   repository.TasksRepo
}

func NewSendDeclarationProcessor(botClient *bot.Bot, tasksService task.TasksService, logger logger.Logger, repo repository.TasksRepo) *SendDeclarationProcessor {
	return &SendDeclarationProcessor{
		bot:    botClient,
		tasksS: tasksService,
		logger: logger,
		repo:   repo,
	}
}

func (p *SendDeclarationProcessor) Process(ctx context.Context, t db.Tasks) error {
	text := fmt.Sprintf("Пожалуйста, отправьте декларацию за отчётный месяц — не забудьте, спасибо!")
	if err := p.bot.SendMessage(ctx, t.EntrepreneurID, text); err != nil {
		p.logger.Errorf("telegram send error for %d: %v", t.EntrepreneurID, err)
		return err
	}

	if err := p.tasksS.CompleteNotification(ctx, t.ID); err != nil {
		p.logger.Info("failed to mark finished task %v: %v", t, err)
		return err
	}

	next := &domain.Task{
		RunAt:          time.Now().Add(24 * time.Hour),
		Status:         "pending",
		Type:           "send_declaration",
		EntrepreneurID: t.EntrepreneurID,
	}

	if err := p.repo.CreateBatch(ctx, []*domain.Task{next}); err != nil {
		p.logger.Errorf("failed to create followup task for %+v: %v", t, err)
		return err
	}

	p.logger.Info("send_declaration processed for %d", t.EntrepreneurID)
	return nil
}
