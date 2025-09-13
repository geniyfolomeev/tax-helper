package db

import (
	"time"

	"github.com/shopspring/decimal"
)

type Entrepreneur struct {
	ID              int64 `gorm:"primaryKey"` // Telegram ID
	Status          string
	RegisteredAt    time.Time
	LastSentAt      time.Time
	YearTotalAmount decimal.Decimal
}

type Income struct {
	ID             int64 `gorm:"primaryKey"`
	EntrepreneurID int64
	Entrepreneur   Entrepreneur `gorm:"foreignKey:EntrepreneurID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Date           time.Time
	Amount         decimal.Decimal
	SourceAmount   decimal.Decimal
	SourceCurrency string
}

type Tasks struct {
	ID             int64 `gorm:"primaryKey"`
	EntrepreneurID int64
	Entrepreneur   Entrepreneur `gorm:"foreignKey:EntrepreneurID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status         string
	Type           string
	RunAt          time.Time
}
