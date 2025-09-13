package db

import (
	"time"

	"github.com/shopspring/decimal"
)

type Entrepreneur struct {
	ID              uint `gorm:"primaryKey"` // Telegram ID
	Status          string
	RegisteredAt    time.Time
	LastSentAt      time.Time
	YearTotalAmount decimal.Decimal `gorm:"type:numeric(18,2)"`
}

type Income struct {
	ID             uint `gorm:"primaryKey"`
	TelegramID     uint
	TelegramUser   Entrepreneur `gorm:"foreignKey:TelegramID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Date           time.Time
	Amount         decimal.Decimal `gorm:"type:numeric(18,2)"`
	SourceAmount   decimal.Decimal `gorm:"type:numeric(18,2)"`
	SourceCurrency string
}

type Tasks struct {
	ID           uint `gorm:"primaryKey"`
	TelegramID   uint
	TelegramUser Entrepreneur `gorm:"foreignKey:TelegramID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status       string
	Type         string
	RunAt        time.Time
}
