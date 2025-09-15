package processors

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/db"
	"tax-helper/internal/infrastructure/repository"
	"tax-helper/internal/logger"
	"tax-helper/internal/service/task"
	"time"
)

type TxManager interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}
type SendDeclarationProcessor struct {
	bot       *tgbotapi.BotAPI
	tasksS    task.TasksService
	logger    logger.Logger
	repo      repository.TasksRepository
	txManager TxManager
}

func NewSendDeclarationProcessor(botClient *tgbotapi.BotAPI, tasksService task.TasksService, logger logger.Logger, repo repository.TasksRepository, txManager TxManager) *SendDeclarationProcessor {
	return &SendDeclarationProcessor{
		bot:       botClient,
		tasksS:    tasksService,
		logger:    logger,
		repo:      repo,
		txManager: txManager,
	}
}

func (p *SendDeclarationProcessor) Process(ctx context.Context, t db.Tasks) error {
	text := "Пожалуйста, отправьте декларацию за отчётный месяц — не забудьте, спасибо!"
	msg := tgbotapi.NewMessage(t.EntrepreneurID, text)
	if _, err := p.bot.Send(msg); err != nil {
		p.logger.Errorf("telegram send error for %d: %v", t.EntrepreneurID, err)
		return err
	}

	err := p.txManager.Transaction(ctx, func(txCtx context.Context) error {
		if err := p.tasksS.CompleteNotification(txCtx, t.ID); err != nil {
			p.logger.Error("failed to mark finished task %v: %v", t, err)
			return err
		}

		next := &domain.Task{
			RunAt:          time.Now().Add(24 * time.Hour),
			Status:         "ready",
			Type:           "submit_declaration",
			EntrepreneurID: t.EntrepreneurID,
		}

		if err := p.repo.CreateBatch(txCtx, []*domain.Task{next}); err != nil {
			p.logger.Errorf("failed to create followup task for %+v: %v", t, err)
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	p.logger.Info("send_declaration processed for %d", t.EntrepreneurID)
	return nil
}
