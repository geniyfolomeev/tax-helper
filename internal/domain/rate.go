package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type Rate struct {
	Code      string
	Quantity  int
	Value     decimal.Decimal
	Date      time.Time
	SrcAmount decimal.Decimal
}

func (r *Rate) ConvertToGel() decimal.Decimal {
	return r.SrcAmount.Mul(r.Value).Div(decimal.NewFromInt(int64(r.Quantity))).Round(2)
}
