package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Income struct {
	EntrepreneurID uint
	Date           time.Time
	Amount         decimal.Decimal
	SourceAmount   decimal.Decimal
	SourceCurrency string
}

func (i *Income) Validate() error {
	return nil
}

func (i *Income) SetAmount(amount decimal.Decimal) {
	i.Amount = amount
}
