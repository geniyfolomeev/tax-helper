package repository

import (
	"errors"
	"fmt"
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
