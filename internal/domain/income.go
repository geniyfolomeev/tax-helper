package domain

import "time"

type MonthIncome struct {
	ID             uint
	Date           time.Time
	Amount         float64
	SourceAmount   float64
	SourceCurrency string
}
