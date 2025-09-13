package rate

import (
	"context"
	"fmt"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/http"
	"tax-helper/internal/logger"
	"time"

	"github.com/shopspring/decimal"
)

type RatesClient interface {
	GetCurrencyRates(ctx context.Context, path string) ([]http.CurrencyRateResponse, error)
}

type CurrencyRateService struct {
	client RatesClient
	logger logger.Logger
}

func NewService(client RatesClient, l logger.Logger) *CurrencyRateService {
	return &CurrencyRateService{client: client, logger: l}
}

func (s *CurrencyRateService) buildPath(currency string, date time.Time) string {
	return fmt.Sprintf("/gw/api/ct/monetarypolicy/currencies/?currencies=%s&date=%s",
		currency,
		date.Format("2006-01-02"),
	)
}

func (s *CurrencyRateService) GetRate(ctx context.Context, amount decimal.Decimal, currency string, date time.Time) (*domain.Rate, error) {
	rates, err := s.client.GetCurrencyRates(ctx, s.buildPath(currency, date))
	if err != nil {
		s.logger.Errorf("failed to get rates: %v", err)
		return nil, err
	}
	if len(rates) == 0 {
		return nil, fmt.Errorf("no rates found for currency %s", currency)
	}

	for _, c := range rates[0].Currencies {
		if c.Code == currency {
			return &domain.Rate{
				Code:      c.Code,
				Quantity:  c.Quantity,
				Value:     decimal.NewFromFloat(c.Rate),
				Date:      c.Date,
				SrcAmount: amount,
			}, nil
		}
	}

	return nil, fmt.Errorf("currency %s not found in rates for %s", currency, date.Format("2006-01-02"))
}
