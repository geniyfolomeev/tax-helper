package repository

import (
	"context"
	"gorm.io/gorm"
	"time"
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

type TaskRepository interface {
	GetPendingTasks(ctx context.Context) ([]Tasks, error)
	MarkAsNotified(ctx context.Context, id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) GetPendingTasks(ctx context.Context) ([]Tasks, error) {
	var tasks []Tasks
	if err := r.db.WithContext(ctx).
		Where("run_at <= ? AND notified = ?", time.Now(), false).
		Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *taskRepository) MarkAsNotified(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&Tasks{}).
		Where("id = ?", id).
		Update("status", "done").Error
}
