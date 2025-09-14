package domain

import "context"

type Notifier interface {
	Process(ctx context.Context, task Task) error
}
