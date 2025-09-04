package domain

import (
	"fmt"
	"time"
)

const (
	sendingDayFrom = 1
	sendingDayTo   = 15
)

// CurrentTimeFn - function to get current time, test purpose only
var currentTimeFn = time.Now

type Entrepreneur struct {
	TelegramID      uint
	Status          string
	RegisteredAt    time.Time // Entrepreneur registered at RegisteredAt, it's not CreatedAt!
	LastSentAt      time.Time // Entrepreneur sent last declaration at LastSentAt (Could be zero date!)
	YearTotalAmount float64
}

func NewEntrepreneur(tgID uint, status string, regAt, lastAt time.Time, yta float64) *Entrepreneur {
	return &Entrepreneur{
		TelegramID:      tgID,
		Status:          status,
		RegisteredAt:    regAt,
		LastSentAt:      lastAt,
		YearTotalAmount: yta,
	}
}

func (e *Entrepreneur) Validate() error {
	if time.Now().Before(e.RegisteredAt) {
		return fmt.Errorf("%w: registration date is in the future", ErrValidation)
	}
	return nil
}

func (e *Entrepreneur) IsActive() bool {
	return e.Status == "active"
}

// Calculate first month when Entrepreneur have to send tax declaration
func (e *Entrepreneur) getFirstSendingInterval() (time.Time, time.Time) {
	nextMonth := e.RegisteredAt.AddDate(0, 1, 0)
	from := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())
	to := time.Date(nextMonth.Year(), nextMonth.Month(), 15, 0, 0, 0, 0, nextMonth.Location())
	return from, to
}

func (e *Entrepreneur) isFirstDeclarationSent() bool {
	from, _ := e.getFirstSendingInterval()
	return !e.LastSentAt.IsZero() && (e.LastSentAt.After(from) || e.LastSentAt.Equal(from)) && e.LastSentAt.Before(currentTimeFn())
}

func (e *Entrepreneur) CalculateNextDeclarationDate() time.Time {
	now := currentTimeFn()
	loc := now.Location()

	if !e.isFirstDeclarationSent() {
		firstFrom, _ := e.getFirstSendingInterval()
		return firstFrom
	}

	curMonthFrom := time.Date(now.Year(), now.Month(), sendingDayFrom, 0, 0, 0, 0, loc)
	nextMonth := now.AddDate(0, 1, 0)
	nextMonthFrom := time.Date(nextMonth.Year(), nextMonth.Month(), sendingDayFrom, 0, 0, 0, 0, loc)

	if now.Day() > sendingDayTo {
		return nextMonthFrom
	}

	if e.LastSentAt.Before(curMonthFrom) {
		return curMonthFrom
	}

	return nextMonthFrom
}
