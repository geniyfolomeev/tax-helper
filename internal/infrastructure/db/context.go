package db

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

func setTx(ctx context.Context, tx *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func getTx(ctx context.Context) *gorm.DB {
	return ctx.Value(txKey{}).(*gorm.DB)
}
