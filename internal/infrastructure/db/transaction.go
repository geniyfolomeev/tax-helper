package db

import (
	"context"

	"gorm.io/gorm"
)

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	transaction := func(tx *gorm.DB) error {
		return fn(setTx(ctx, tx))
	}
	return m.db.Session(&gorm.Session{
		SkipDefaultTransaction: true,
	}).Transaction(transaction)
}
