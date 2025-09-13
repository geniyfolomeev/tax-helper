package domain

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

const (
	sendingDayFrom = 1
	sendingDayTo   = 15
)

type Entrepreneur struct {
	TelegramID      uint
	Status          string
	RegisteredAt    time.Time // Entrepreneur registered at RegisteredAt, it's not CreatedAt!
	LastSentAt      time.Time // Entrepreneur sent last declaration at LastSentAt (Could be zero date!)
	YearTotalAmount decimal.Decimal
}

func (e *Entrepreneur) Validate() error {
	if e.YearTotalAmount.IsNegative() {
		return fmt.Errorf("%w: yearly amount cannot be negative", ErrValidation)
	}
	if currentTimeFn().Before(e.RegisteredAt) {
		return fmt.Errorf("%w: registration date cannot be in the future", ErrValidation)
	}
	if e.LastSentAt.IsZero() && !e.YearTotalAmount.IsZero() {
		return fmt.Errorf("%w: yearly amount must be zero until the first declaration is sent", ErrValidation)
	}
	if !e.LastSentAt.IsZero() && e.YearTotalAmount.IsZero() {
		return fmt.Errorf("%w: yearly amount cannot be zero after at least one declaration has been sent", ErrValidation)
	}
	if !e.LastSentAt.IsZero() && e.RegisteredAt.After(e.LastSentAt) {
		return fmt.Errorf("%w: registration date cannot be after your last declaration", ErrValidation)
	}
	if !e.LastSentAt.IsZero() && e.LastSentAt.After(currentTimeFn()) {
		return fmt.Errorf("%w: last declaration date cannot be in the future", ErrValidation)
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
