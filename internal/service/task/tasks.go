package task

import (
	"context"
	"tax-helper/internal/infrastructure/db"
	"tax-helper/internal/infrastructure/repository"
	"time"
)

type TasksService interface {
	GetDueTasks(ctx context.Context, now time.Time) ([]db.Tasks, error)
	CompleteNotification(ctx context.Context, id int64) error
}

type tasksService struct {
	repo repository.TasksRepository
}

func NewService(repo repository.TasksRepository) TasksService {
	return &tasksService{repo: repo}
}

func (s *tasksService) GetDueTasks(ctx context.Context, timeNow time.Time) ([]db.Tasks, error) {
	return s.repo.GetReadyTasks(ctx, timeNow)
}

func (s *tasksService) CompleteNotification(ctx context.Context, id int64) error {
	return s.repo.MarkAsNotified(ctx, id)
}
