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
	"tax-helper/internal/infrastructure/bot"
	"tax-helper/internal/infrastructure/db"
	"tax-helper/internal/infrastructure/http"
	"tax-helper/internal/infrastructure/processors"
	"tax-helper/internal/infrastructure/repository"
	logging "tax-helper/internal/logger"
	"tax-helper/internal/scheduler"
	"tax-helper/internal/service/entrepreneur"
	"tax-helper/internal/service/income"
	"tax-helper/internal/service/rate"
	"tax-helper/internal/service/task"
	"time"
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
	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("failed to init bot: %v", err)
	}

	txManager := db.NewTxManager(database.DefaultConnection())
	httpClient := http.NewClient(cfg.RateUrl, time.Second*5, logger)

	entrepreneurRepo := repository.NewEntrepreneurRepo(database)
	tasksRepo := repository.NewTasksRepo(database)
	incomeRepo := repository.NewIncomeRepo(database)

	rateService := rate.NewService(httpClient, logger)
	entrepreneurService := entrepreneur.NewService(entrepreneurRepo, tasksRepo, txManager)
	incomeService := income.NewIncomeService(incomeRepo, rateService, logger)
	tasksService := task.NewService(tasksRepo)

	tgBot, err := bot.NewBot(entrepreneurService, incomeService, logger, botApi)
	processors.CompositionsAdapters(tgBot)
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
