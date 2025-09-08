package income

import (
	"context"
	"errors"
	"fmt"
	"tax-helper/internal/domain"
	"tax-helper/internal/logger"
	"time"

	"github.com/shopspring/decimal"
)

type UserIncomeRepo interface {
	AddIncome(ctx context.Context, income *domain.Income) error
	GetIncomeByPeriod(ctx context.Context, tgID uint, dateFrom time.Time, dateTo time.Time) ([]*domain.Income, error)
}

type RateService interface {
	GetRate(ctx context.Context, amount decimal.Decimal, currency string, date time.Time) (*domain.Rate, error)
}

type Service struct {
	incomeRepo  UserIncomeRepo
	rateService RateService
	logger      logger.Logger
}

func NewIncomeService(ir UserIncomeRepo, rs RateService, l logger.Logger) *Service {
	return &Service{
		incomeRepo:  ir,
		rateService: rs,
		logger:      l,
	}
}

func (s *Service) AddIncome(ctx context.Context, tgID uint, amount float64, currency string, date time.Time) error {
	income := &domain.Income{
		EntrepreneurID: tgID,
		Amount:         decimal.NewFromFloat(0),
		SourceAmount:   decimal.NewFromFloat(amount),
		SourceCurrency: currency,
		Date:           date,
	}
	err := income.Validate()
	if err != nil {
		return err
	}

	currencyRate, err := s.rateService.GetRate(ctx, income.SourceAmount, currency, date)
	if err != nil {
		return err
	}

	income.SetAmount(currencyRate.ConvertToGel())
	err = s.incomeRepo.AddIncome(ctx, income)
	if err != nil {
		if errors.Is(err, domain.ErrEntrepreneurNotFound) {
			return fmt.Errorf("%w: you need to register first", err)
		}
		return err
	}
	return nil
}

// GetActualIncome - get Entrepreneur income for previous and current month
func (s *Service) GetActualIncome(ctx context.Context, tgID uint) (*domain.ActualIncome, error) {
	prevFrom, _, _, curTo := domain.GetActualTimeIntervals()
	incomes, err := s.incomeRepo.GetIncomeByPeriod(ctx, tgID, prevFrom, curTo)
	if err != nil {
		return nil, err
	}
	actualIncome := domain.NewActualIncome(incomes)
	return actualIncome, nil
}
