package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tax-helper/internal/config"
	"tax-helper/internal/infrastructure/bot"
	"tax-helper/internal/infrastructure/db"
	"tax-helper/internal/infrastructure/repository"
	"tax-helper/internal/logger"
	"tax-helper/internal/scheduler"
	"tax-helper/internal/service"
	"time"
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

	txManager := db.NewTxManager(database.DefaultConnection())

	entrepreneurRepo := repository.NewEntrepreneurRepo(database)
	tasksRepo := repository.NewTasksRepo(database)

	taxService := service.NewTaxService(entrepreneurRepo, tasksRepo, txManager)
	tasksService := service.NewTaskService(tasksRepo)

	tgBot, err := bot.NewBot(cfg, taxService)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err = tgBot.Run(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	s := scheduler.NewScheduler(tasksService, time.Hour, tgBot)
	s.Start(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sig:
		logger.Info("shutting down...")
		cancel()
	case <-ctx.Done():
	}

	wg.Wait()
}
