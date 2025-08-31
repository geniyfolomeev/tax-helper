package main

import (
	"log"
	"tax-helper/config"
	"tax-helper/internal/infrastructure/db"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	database, err := db.NewDB(cfg)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: проверять нужны ли вообще миграции или схема БД актуальная
	err = database.Migrate()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migration completed")
}
