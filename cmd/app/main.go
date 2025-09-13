package main

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tax-helper/internal/config"
	"tax-helper/internal/infrastructure/adapters"
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
	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("failed to init bot: %v", err)
	}

	entrepreneurRepo := repository.NewEntrepreneurRepo(database.Conn)
	entrepreneurTasksRepo := repository.NewEntrepreneurTasksRepo(database.Conn)
	tasksRepo := repository.NewTaskRepository(database.Conn)
	entrepreneurService := service.NewEntrepreneurService(entrepreneurRepo, entrepreneurTasksRepo)
	notifiers := adapters.CompositionsAdapters(botApi)
	tasksService := service.NewTaskService(tasksRepo, notifiers)

	tgBot, err := bot.NewBot(cfg, entrepreneurService, botApi)
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

	s := scheduler.NewScheduler(tasksService, time.Hour)
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
