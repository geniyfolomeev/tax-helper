package repository

import (
	"tax-helper/internal/domain"
	"time"

	"gorm.io/gorm"
)

type Tasks struct {
	ID         uint `gorm:"primaryKey"`
	TelegramID uint
	Status     string
	Type       string
	RunAt      time.Time
}

type TasksRepo struct {
	db *gorm.DB
}

func NewTasksRepo(db *gorm.DB) *TasksRepo {
	return &TasksRepo{db: db}
}

func (r *TasksRepo) CreateBatch(tasks []*domain.Task) error {
	models := make([]*Tasks, len(tasks))
	for i, t := range tasks {
		models[i] = &Tasks{
			TelegramID: t.TelegramID,
			Status:     t.Status,
			Type:       t.Type,
			RunAt:      t.RunAt,
		}
	}
	return r.db.Create(&models).Error
}
