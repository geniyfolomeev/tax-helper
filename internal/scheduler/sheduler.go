package scheduler

import (
	"context"
	"fmt"
	"tax-helper/internal/domain"
	"tax-helper/internal/service/task"
	"time"
)

type Scheduler struct {
	service   *task.TasksService
	ticker    *time.Ticker
	processor map[string]domain.Notifier
}

func NewScheduler(service *task.TasksService, interval time.Duration) *Scheduler {
	return &Scheduler{
		service: service,
		ticker:  time.NewTicker(interval),
	}
}
func (s *Scheduler) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				s.ticker.Stop()
				return
			case <-s.ticker.C:
				tasks, err := s.service.GetDueTasks(ctx)
				if err != nil {
					fmt.Printf("Ошибка при получении задач: %v\n", err)
					continue
				}

				for _, task := range tasks {
					if processor, ok := s.processor[task.Type]; ok {
						if err := processor.Process(ctx, task); err != nil {
							fmt.Printf("Ошибка при обработке задачи %d: %v\n", task.TelegramID, err)
						}
					} else {
						fmt.Printf("Нет процессора для задачи типа %s\n", task.Type)
					}
				}
			}
		}
	}()
}
