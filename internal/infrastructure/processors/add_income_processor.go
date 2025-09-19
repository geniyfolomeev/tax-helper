package processors

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/db"
	"tax-helper/internal/infrastructure/repository"
	"tax-helper/internal/logger"
	"time"
)

type AddIncomeProcessor struct {
	bot       *tgbotapi.BotAPI
	logger    logger.Logger
	repo      repository.TasksRepository
	txManager TxManager
}

func NewAddIncomeProcessor(botClient *tgbotapi.BotAPI, logger logger.Logger, repo repository.TasksRepository, txManager TxManager) *AddIncomeProcessor {
	return &AddIncomeProcessor{bot: botClient, logger: logger, repo: repo, txManager: txManager}
}

func (p *AddIncomeProcessor) Process(ctx context.Context, t db.Tasks) error {
	text := "Пора внести сведения о доходе за отчётный месяц — отправьте, пожалуйста."
	msg := tgbotapi.NewMessage(t.EntrepreneurID, text)
	if _, err := p.bot.Send(msg); err != nil {
		p.logger.Errorf("telegram send error for %d: %v", t.EntrepreneurID, err)
		return err
	}

	err := p.txManager.Transaction(ctx, func(txCtx context.Context) error {
		if err := p.repo.MarkAsNotified(txCtx, t.ID); err != nil {
			p.logger.Info("failed to mark finished task %v: %v", t, err)
			return err
		}

		next := &domain.Task{
			RunAt:          time.Now().Add(24 * time.Hour),
			Status:         "ready",
			Type:           "add_income",
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

	p.logger.Info("add_income processed for %d", t.EntrepreneurID)
	return nil
}
