package domain

import (
	"context"
	"tax-helper/internal/infrastructure/db"
)

type Notifier interface {
	Process(ctx context.Context, task db.Tasks) error
}
