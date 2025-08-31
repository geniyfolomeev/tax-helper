package main

import (
	"log"
	"tax-helper/config"
	"tax-helper/internal/infrastructure/db"
	"tax-helper/internal/infrastructure/http"
	"tax-helper/internal/infrastructure/repository"
	"tax-helper/internal/service"
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
	entrepreneurRepo := repository.NewEntrepreneurRepo(database.Conn)
	entrepreneurService := service.NewEntrepreneurService(entrepreneurRepo)
	srv := http.New(entrepreneurService)
	if err = srv.Run(cfg.HTTPPort); err != nil {
		log.Fatal(err)
	}
}
