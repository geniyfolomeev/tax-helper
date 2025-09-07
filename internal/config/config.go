package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDsn                string
	DBMaxOpenConnections int
	DBMaxIdleConnections int

	BotToken string

	LoggingMode string
}

func NewConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
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
		LoggingMode:          os.Getenv("LOGGING_MODE"),
	}, nil
}
