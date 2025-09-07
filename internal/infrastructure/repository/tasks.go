package repository

import (
	"context"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/db"
)

type TasksRepo struct {
	db *db.DB
}

func NewTasksRepo(db *db.DB) *TasksRepo {
	return &TasksRepo{db: db}
}

func (r *TasksRepo) CreateBatch(ctx context.Context, tasks []*domain.Task) error {
	models := make([]*db.Tasks, len(tasks))
	for i, t := range tasks {
		models[i] = &db.Tasks{
			TelegramID: t.TelegramID,
			Status:     t.Status,
			Type:       t.Type,
			RunAt:      t.RunAt,
		}
	}
	return r.db.Connection(ctx).Create(&models).Error
}
