package config

import (
	"os"
	"strconv"
	"tax-helper/internal/logger"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDsn                string
	DBMaxOpenConnections int
	DBMaxIdleConnections int

	BotToken string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found, using system env")
	}
	maxOpenConnections, err := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTIONS"))
	if err != nil {
		return nil, err
	}
	maxIdleConnections, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	if err != nil {
		return nil, err
	}
	return &Config{
		DBDsn:                os.Getenv("DB_DSN"),
		DBMaxOpenConnections: maxOpenConnections,
		DBMaxIdleConnections: maxIdleConnections,
		BotToken:             os.Getenv("BOT_TOKEN"),
	}, nil
}
