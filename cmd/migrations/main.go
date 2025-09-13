package main

import (
	"log"
	"tax-helper/internal/config"
	"tax-helper/internal/infrastructure/db"
	logging "tax-helper/internal/logger"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, err := logging.NewLogger(cfg.LoggingMode)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	database, err := db.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: проверять нужны ли вообще миграции или схема БД актуальная
	err = database.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	logger.Info("Migrations completed")
}
