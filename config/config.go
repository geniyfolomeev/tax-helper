package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDsn                string
	DBMaxOpenConnections int
	DBMaxIdleConnections int

	HTTPPort string
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
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
		HTTPPort:             os.Getenv("HTTP_PORT"),
	}, nil
}
