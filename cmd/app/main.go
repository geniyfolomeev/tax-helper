package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tax-helper/config"
	"tax-helper/internal/infrastructure/bot"
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

	srv := http.NewServer(cfg, entrepreneurService)
	tgBot, err := bot.NewBot(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = srv.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = tgBot.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sig:
		log.Println("shutting down...")
		cancel()
	case <-ctx.Done():
	}

	wg.Wait()
}
