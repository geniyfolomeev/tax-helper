package repository

import (
	"errors"
	"fmt"
	"strings"
	"tax-helper/internal/domain"

	"gorm.io/gorm"
)

type EntrepreneurTasksRepo struct {
	db *gorm.DB
}

func NewEntrepreneurTasksRepo(db *gorm.DB) *EntrepreneurTasksRepo {
	return &EntrepreneurTasksRepo{
		db: db,
	}
}

func (r *EntrepreneurTasksRepo) Create(
	e *domain.Entrepreneur,
	tasks []*domain.Task,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		model := &Entrepreneur{
			ID:              e.TelegramID,
			Status:          e.Status,
			RegisteredAt:    e.RegisteredAt,
			LastSentAt:      e.LastSentAt,
			YearTotalAmount: e.YearTotalAmount,
		}
		err := tx.Create(model).Error
		if err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate") {
				return fmt.Errorf("%w: telegram_id = %d", domain.ErrEntrepreneurAlreadyExists, e.TelegramID)
			}
			return err
		}
		models := make([]*Tasks, len(tasks))
		for i, t := range tasks {
			models[i] = &Tasks{
				TelegramID: t.TelegramID,
				Status:     t.Status,
				Type:       t.Type,
				RunAt:      t.RunAt,
			}
		}
		return tx.Create(&models).Error
	})
}
