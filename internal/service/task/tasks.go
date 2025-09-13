package task

import (
	"context"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/db"
	"tax-helper/internal/infrastructure/repository"
)

type TasksService struct {
	repo     *repository.TasksRepo
	notifier map[string]domain.Notifier
}

func NewService(repo *repository.TasksRepo, notifier map[string]domain.Notifier) *TasksService {
	return &TasksService{repo: repo,
		notifier: notifier}
}

func (s *TasksService) GetDueTasks(ctx context.Context) ([]db.Tasks, error) {
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
	for _, r := range reminders {
		// тут ищем нужный адаптер
		adapter, ok := s.notifier["tg"]
		if !ok {
			continue // нет подходящего адаптера
		}
		if err := adapter.SendMessage(int64(r.TelegramID), "test"); err != nil {
			return err
		}
		if err := s.repo.MarkAsNotified(ctx, r.ID); err != nil {
			return err
		}
	}
	return nil
}
