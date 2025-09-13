package repository

import (
	"context"
	"strings"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/db"
	"time"

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

func (r *IncomeRepo) GetIncomeByPeriod(ctx context.Context, tgID uint, dateFrom time.Time, dateTo time.Time) ([]*domain.Income, error) {
	var rows []*db.Income
	err := r.db.Connection(ctx).
		Where("telegram_id = ? AND date BETWEEN ? AND ?", tgID, dateFrom, dateTo).
		Find(&rows).Error

	if err != nil {
		return nil, err
	}

	incomes := make([]*domain.Income, len(rows))
	for i, row := range rows {
		incomes[i] = &domain.Income{
			EntrepreneurID: row.TelegramID,
			Date:           row.Date,
			Amount:         row.Amount,
			SourceAmount:   row.SourceAmount,
			SourceCurrency: row.SourceCurrency,
		}
	}
	return incomes, nil
}
