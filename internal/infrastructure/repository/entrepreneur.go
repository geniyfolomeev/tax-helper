package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/db"

	"gorm.io/gorm"
)

type EntrepreneurRepo struct {
	db *db.DB
}

func NewEntrepreneurRepo(db *db.DB) *EntrepreneurRepo {
	return &EntrepreneurRepo{db: db}
}

func (r *EntrepreneurRepo) GetByID(ctx context.Context, id uint) (*domain.Entrepreneur, error) {
	var e db.Entrepreneur
	if err := r.db.Connection(ctx).Where(&db.Entrepreneur{ID: id}).First(&e).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: id = %d", domain.ErrEntrepreneurNotFound, id)
		}
		return nil, err
	}

	return &domain.Entrepreneur{
		TelegramID:      e.ID,
		Status:          e.Status,
		RegisteredAt:    e.RegisteredAt,
		YearTotalAmount: e.YearTotalAmount,
		LastSentAt:      e.LastSentAt,
	}, nil
}

func (r *EntrepreneurRepo) Create(ctx context.Context, e *domain.Entrepreneur) error {
	model := &db.Entrepreneur{
		ID:              e.TelegramID,
		Status:          e.Status,
		RegisteredAt:    e.RegisteredAt,
		LastSentAt:      e.LastSentAt,
		YearTotalAmount: e.YearTotalAmount,
	}
	err := r.db.Connection(ctx).Create(model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("%w: telegram_id = %d", domain.ErrEntrepreneurAlreadyExists, e.TelegramID)
		}
		return err
	}
	return nil
}
