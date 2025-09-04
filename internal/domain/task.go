package domain

import "time"

type Task struct {
	TelegramID uint
	Status     string
	Type       string
	RunAt      time.Time
}

func NewTask(tgID uint, status, taskType string, runAt time.Time) *Task {
	return &Task{
		TelegramID: tgID,
		Status:     status,
		Type:       taskType,
		RunAt:      runAt,
	}
}
