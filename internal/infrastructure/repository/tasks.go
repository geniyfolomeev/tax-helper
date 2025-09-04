package repository

import (
	"time"

	"gorm.io/gorm"
)

type Tasks struct {
	ID         uint `gorm:"primaryKey"`
	TelegramID uint
	Status     string
	Type       string
	RunAt      time.Time
}

type TasksRepo struct {
	db *gorm.DB
}

func NewTasksRepo(db *gorm.DB) *TasksRepo {
	return &TasksRepo{db: db}
}
