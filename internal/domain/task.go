package domain

import "time"

type Task struct {
	EntrepreneurID int64
	Status         string
	Type           string
	RunAt          time.Time
}

func NewTask(tgID int64, status, taskType string, runAt time.Time) *Task {
	return &Task{
		EntrepreneurID: tgID,
		Status:         status,
		Type:           taskType,
		RunAt:          runAt,
	}
}
