package db

import (
	"tax-helper/internal/config"
	"tax-helper/internal/infrastructure/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Conn *gorm.DB
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
	return &DB{Conn: conn}, nil
}

func (db *DB) Migrate() error {
	return db.Conn.AutoMigrate(&repository.Entrepreneur{}, &repository.Tasks{})
}
