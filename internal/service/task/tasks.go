package task

import (
	"context"
	"tax-helper/internal/infrastructure/db"
	"tax-helper/internal/infrastructure/repository"
)

type TasksService struct {
	repo *repository.TasksRepo
}

func NewService(repo *repository.TasksRepo) *TasksService {
	return &TasksService{repo: repo}
}

func (s *TasksService) GetDueTasks(ctx context.Context) ([]db.Tasks, error) {
	return s.repo.GetPendingTasks(ctx)
}

func (s *TasksService) CompleteNotification(ctx context.Context, id int64) error {
	return s.repo.MarkAsNotified(ctx, id)
}
