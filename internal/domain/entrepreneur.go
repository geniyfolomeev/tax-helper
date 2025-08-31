package domain

import (
	"fmt"
	"time"
)

type Entrepreneur struct {
	TelegramID       uint
	Status           string
	RegistrationDate time.Time
	YearTotalAmount  float64
}

func (e *Entrepreneur) Validate() error {
	if time.Now().Before(e.RegistrationDate) {
		return fmt.Errorf("%w: registration date is in the future", ErrValidation)
	}
	return nil
}
