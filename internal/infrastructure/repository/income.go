package repository

import (
	"context"
	"strings"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/db"

	"gorm.io/gorm"
)

type IncomeRepo struct {
	db *db.DB
}

func NewIncomeRepo(db *db.DB) *IncomeRepo {
	return &IncomeRepo{db: db}
}

func (r *IncomeRepo) AddIncome(ctx context.Context, income *domain.Income) error {
	model := &db.Income{
		TelegramID:     income.EntrepreneurID,
		Date:           income.Date,
		Amount:         income.Amount,
		SourceAmount:   income.SourceAmount,
		SourceCurrency: income.SourceCurrency,
	}
	err := r.db.Connection(ctx).Create(model).Error
	if err != nil && strings.Contains(err.Error(), gorm.ErrForeignKeyViolated.Error()) {
		return domain.ErrEntrepreneurNotFound
	}
	return nil
}
