package service

import (
	"context"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/repository"
	"tax-helper/internal/logger"
)

type TaskService interface {
	GetDueTasks(ctx context.Context) ([]repository.Tasks, error)
	CompleteNotification(ctx context.Context, id uint) error
	ProcessDueReminders(ctx context.Context) error
}

type TasksService struct {
	repo     repository.TaskRepository
	notifier map[string]domain.Notifier
}

func NewTaskService(repo repository.TaskRepository, notifier map[string]domain.Notifier) TaskService {
	return &TasksService{repo: repo,
		notifier: notifier}
}

func (s *TasksService) GetDueTasks(ctx context.Context) ([]repository.Tasks, error) {
	return s.repo.GetPendingTasks(ctx)
}

func (s *TasksService) CompleteNotification(ctx context.Context, id uint) error {
	return s.repo.MarkAsNotified(ctx, id)
}

func (s *TasksService) ProcessDueReminders(ctx context.Context) error {
	reminders, err := s.repo.GetPendingTasks(ctx)
	if err != nil {
		return err
	}
	logger.Info("tasks", reminders)
	for _, r := range reminders {
		adapter, ok := s.notifier[r.Type]
		if !ok {
			logger.Error("no adapters", nil)
			continue // нет подходящего адаптера
		}
		if err := adapter.SendMessage(int64(r.TelegramID), "asd"); err != nil {
			return err
		}
		if err := s.repo.MarkAsNotified(ctx, r.ID); err != nil {
			return err
		}
	}
	return nil
}
