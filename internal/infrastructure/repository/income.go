package repository

import (
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/db"
)

type IncomeRepo struct {
	db *db.DB
}

func NewIncomeRepo(db *db.DB) *IncomeRepo {
	return &IncomeRepo{db: db}
}

func (r *IncomeRepo) GetByID(id uint) (*domain.MonthIncome, error) {
	return nil, nil
}
