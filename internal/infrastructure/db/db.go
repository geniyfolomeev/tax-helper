package db

import (
	"context"
	"tax-helper/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	conn *gorm.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	conn, err := gorm.Open(postgres.Open(cfg.DBDsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	dbSettings, err := conn.DB()
	if err != nil {
		return nil, err
	}
	dbSettings.SetMaxOpenConns(cfg.DBMaxOpenConnections)
	dbSettings.SetMaxIdleConns(cfg.DBMaxIdleConnections)
	return &DB{conn: conn}, nil
}

func (db *DB) DefaultConnection() *gorm.DB {
	return db.conn
}

func (db *DB) Connection(ctx context.Context) *gorm.DB {
	txCtx, ok := getTx(ctx)
	if ok {
		return txCtx
	}
	return db.DefaultConnection()
}

func (db *DB) Migrate() error {
	return db.conn.AutoMigrate(&Entrepreneur{}, &Tasks{}, &Income{})
}
