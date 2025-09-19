package repository

import (
	"context"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/db"
	"time"
)

type TasksRepository interface {
	CreateBatch(ctx context.Context, tasks []*domain.Task) error
	GetReadyTasks(ctx context.Context, timeNow time.Time) ([]db.Tasks, error)
	MarkAsNotified(ctx context.Context, id int64) error
}

type tasksRepo struct {
	db *db.DB
}

func NewTasksRepo(db *db.DB) TasksRepository {
	return &tasksRepo{db: db}
}

func (r *tasksRepo) CreateBatch(ctx context.Context, tasks []*domain.Task) error {
	models := make([]*db.Tasks, len(tasks))
	for i, t := range tasks {
		models[i] = &db.Tasks{
			EntrepreneurID: t.EntrepreneurID,
			Status:         t.Status,
			Type:           t.Type,
			RunAt:          t.RunAt,
		}
	}
	return r.db.Connection(ctx).Create(&models).Error
}

func (r *tasksRepo) GetReadyTasks(ctx context.Context, timeNow time.Time) ([]db.Tasks, error) {
	var dbTasks []db.Tasks

	if err := r.db.Connection(ctx).
		Where("run_at >= ? AND status = ?", timeNow, "ready"). //Todo run_at >= ? или <=
		Find(&dbTasks).Error; err != nil {
		return nil, err
	}
	return dbTasks, nil
}

func (r *tasksRepo) MarkAsNotified(ctx context.Context, id int64) error {
	return r.db.Connection(ctx).Model(&db.Tasks{}).
		Where("id = ?", id).
		Update("status", "done").Error
}
