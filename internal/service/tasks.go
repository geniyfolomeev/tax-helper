package service

import (
	"context"
	"tax-helper/internal/infrastructure/repository"
)

type TaskService interface {
	GetDueTasks(ctx context.Context) ([]repository.Tasks, error)
	CompleteNotification(ctx context.Context, id uint) error
}

type TasksService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &TasksService{repo: repo}
}

func (s *TasksService) GetDueTasks(ctx context.Context) ([]repository.Tasks, error) {
	return s.repo.GetPendingTasks(ctx)
}

func (s *TasksService) CompleteNotification(ctx context.Context, id uint) error {
	return s.repo.MarkAsNotified(ctx, id)
}
