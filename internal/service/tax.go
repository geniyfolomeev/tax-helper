package service

import (
	"context"
	"tax-helper/internal/domain"
	"time"
)

type EntrepreneurRepo interface {
	GetByID(ctx context.Context, id uint) (*domain.Entrepreneur, error)
	Create(ctx context.Context, e *domain.Entrepreneur) error
}

type TasksRepo interface {
	CreateBatch(ctx context.Context, tasks []*domain.Task) error
}

type TxManager interface {
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type TaxService struct {
	entrepreneurRepo EntrepreneurRepo
	tasksRepo        TasksRepo
	txManager        TxManager
}

func NewTaxService(
	er EntrepreneurRepo,
	tr TasksRepo,
	txManager TxManager,
) *TaxService {
	return &TaxService{
		entrepreneurRepo: er,
		tasksRepo:        tr,
		txManager:        txManager,
	}
}

func (s *TaxService) CreateEntrepreneur(ctx context.Context, tgID uint, regAt, lastAt time.Time, yta float64) error {
	e := &domain.Entrepreneur{
		TelegramID:      tgID,
		Status:          "active",
		RegisteredAt:    regAt,
		LastSentAt:      lastAt,
		YearTotalAmount: yta,
	}
	err := e.Validate()
	if err != nil {
		return err
	}

	nextDeclarationDate := e.CalculateNextDeclarationDate()
	tasks := []*domain.Task{
		domain.NewTask(e.TelegramID, "ready", "add_income", nextDeclarationDate),
		domain.NewTask(e.TelegramID, "ready", "submit declaration", nextDeclarationDate),
	}

	return s.txManager.Transaction(ctx, func(ctx context.Context) error {
		err = s.entrepreneurRepo.Create(ctx, e)
		if err != nil {
			return err
		}
		return s.tasksRepo.CreateBatch(ctx, tasks)
	})
}
