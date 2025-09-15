package repository

import (
	"context"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/db"
	"time"
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

func (r *TasksRepo) GetPendingTasks(ctx context.Context) ([]db.Tasks, error) {
	var dbTasks []db.Tasks
	if err := r.db.Connection(ctx).
		Where("run_at <= ? AND status = ?", time.Now(), "active").
		Find(&dbTasks).Error; err != nil {
		return nil, err
	}

	tasks := make([]db.Tasks, 0, len(dbTasks))
	for _, dbTask := range dbTasks {
		tasks = append(tasks, db.Tasks{
			ID:           dbTask.ID,
			TelegramID:   dbTask.TelegramID,
			TelegramUser: dbTask.TelegramUser,
			Status:       dbTask.Status,
			Type:         dbTask.Type,
			RunAt:        dbTask.RunAt,
		})
	}

	return tasks, nil
}

func (r *TasksRepo) MarkAsNotified(ctx context.Context, id uint) error {
	return r.db.Connection(ctx).Model(&db.Tasks{}).
		Where("id = ?", id).
		Update("status", "done").Error
}
