package processors

import (
	"context"
	"fmt"
	"tax-helper/internal/domain"
)

type AddIncomeProcessor struct{}

func NewEmailNotifier() *AddIncomeProcessor {
	return &AddIncomeProcessor{}
}

func (p *AddIncomeProcessor) Process(ctx context.Context, task domain.Task) error {

	fmt.Printf("[AddIncomeProcessor] Добавлен доход: %.2f\n", task.Type)
	return nil
}
