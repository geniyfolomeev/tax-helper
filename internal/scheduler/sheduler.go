package scheduler

import (
	"context"
	"fmt"
	"tax-helper/internal/infrastructure/bot"
	"tax-helper/internal/logger"
	"tax-helper/internal/service"
	"time"
)

type Scheduler struct {
	service  *service.TasksService
	interval time.Duration
	bot      *bot.Bot
}

func NewScheduler(service *service.TasksService, interval time.Duration, bot *bot.Bot) *Scheduler {
	return &Scheduler{
		service:  service,
		interval: interval,
		bot:      bot,
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	ticker := time.NewTicker(s.interval)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				// для дебага
				fmt.Println("Scheduler остановлен")
				return
			case t := <-ticker.C:
				fmt.Println("Проверка напоминаний в:", t)

				tasks, err := s.service.GetDueTasks(ctx)
				if err != nil {
					logger.Error("Ошибка при пометке:", err)
					continue
				}

				if len(tasks) == 0 {
					continue
				}

				for _, task := range tasks {
					if err := s.bot.SendMessage(task.TelegramID, "Заглушка по отпрравке деклеарации"); err != nil {
						logger.Error("Ошибка при пометке:", err)
						continue
					}

					if err := s.service.CompleteNotification(ctx, task.ID); err != nil {
						logger.Error("Ошибка при пометке:", err)
					}
				}
			}
		}
	}()
}
