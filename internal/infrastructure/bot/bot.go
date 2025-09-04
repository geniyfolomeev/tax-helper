package bot

import (
	"context"
	"log"
	"tax-helper/internal/config"
	"tax-helper/internal/infrastructure/bot/commands"
	"tax-helper/internal/logger"
	"tax-helper/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handler interface {
	Command() tgbotapi.BotCommand
	Handle(api *tgbotapi.BotAPI, msg *tgbotapi.Message) (tgbotapi.Message, error)
}

type Bot struct {
	api          *tgbotapi.BotAPI
	updates      tgbotapi.UpdatesChannel
	cmdToHandler map[string]handler
}

func NewBot(cfg *config.Config, e *service.EntrepreneurService) (*Bot, error) {
	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	handlers := []handler{
		commands.NewStartHandler(),
		commands.NewHelpHandler(),
		commands.NewRegisterHandler(e),
	}
	cmdToHandler := map[string]handler{}
	cfgCommands := make([]tgbotapi.BotCommand, 0, len(handlers))
	for _, h := range handlers {
		cfgCommands = append(cfgCommands, h.Command())
		cmdToHandler[h.Command().Command] = h
	}

	botCfg := tgbotapi.SetMyCommandsConfig{
		Commands: cfgCommands,
	}
	_, err = botApi.Request(botCfg)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:          botApi,
		updates:      botApi.GetUpdatesChan(u),
		cmdToHandler: cmdToHandler,
	}, nil
}

func (bot *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	h, ok := bot.cmdToHandler[update.Message.Command()]
	if !ok {
		// TODO: default handler
		return
	}
	_, err := h.Handle(bot.api, update.Message)
	if err != nil {
		log.Println(err)
	}
}

func (bot *Bot) Run(ctx context.Context) error {
	logger.Info("bot started")
	for {
		select {
		case <-ctx.Done():
			logger.Info("bot stopped")
			return nil
		case update := <-bot.updates:
			bot.handleUpdate(update)
		}
	}
}
