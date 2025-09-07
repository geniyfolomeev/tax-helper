package db

import "time"

type Entrepreneur struct {
	ID              uint `gorm:"primaryKey"` // Telegram ID
	Status          string
	RegisteredAt    time.Time
	LastSentAt      time.Time
	YearTotalAmount float64
}

type Income struct {
	ID             uint `gorm:"primaryKey"`
	Date           time.Time
	Amount         float64
	SourceAmount   float64
	SourceCurrency string
}

type Tasks struct {
	ID         uint `gorm:"primaryKey"`
	TelegramID uint
	Status     string
	Type       string
	RunAt      time.Time
}
