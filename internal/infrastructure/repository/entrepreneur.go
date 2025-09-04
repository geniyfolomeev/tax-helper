package repository

import (
	"errors"
	"fmt"
	"strings"
	"tax-helper/internal/domain"
	"time"

	"gorm.io/gorm"
)

type Entrepreneur struct {
	ID              uint `gorm:"primaryKey"` // Telegram ID
	Status          string
	RegisteredAt    time.Time
	LastSentAt      time.Time
	YearTotalAmount float64
}

type EntrepreneurRepo struct {
	db *gorm.DB
}

func NewEntrepreneurRepo(db *gorm.DB) *EntrepreneurRepo {
	return &EntrepreneurRepo{db: db}
}

func (r *EntrepreneurRepo) Create(e *domain.Entrepreneur) error {
	model := &Entrepreneur{
		ID:              e.TelegramID,
		Status:          e.Status,
		RegisteredAt:    e.RegisteredAt,
		LastSentAt:      e.LastSentAt,
		YearTotalAmount: e.YearTotalAmount,
	}
	err := r.db.Create(model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate") {
			return fmt.Errorf("%w: telegram_id = %d", domain.ErrEntrepreneurAlreadyExists, e.TelegramID)
		}
		return err
	}
	return nil
}

func (r *EntrepreneurRepo) GetByID(id uint) (*domain.Entrepreneur, error) {
	var e Entrepreneur
	if err := r.db.Where(&Entrepreneur{ID: id}).First(&e).Error; err != nil {
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
