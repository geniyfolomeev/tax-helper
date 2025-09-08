package bot

import (
	"context"
	"runtime/debug"
	"tax-helper/internal/config"
	"tax-helper/internal/infrastructure/bot/commands"
	"tax-helper/internal/logger"
	"tax-helper/internal/service/entrepreneur"
	"tax-helper/internal/service/income"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type handler interface {
	Command() tgbotapi.BotCommand
	Handle(ctx context.Context, api *tgbotapi.BotAPI, msg *tgbotapi.Message) (tgbotapi.Message, error)
}

type Bot struct {
	api          *tgbotapi.BotAPI
	updates      tgbotapi.UpdatesChannel
	cmdToHandler map[string]handler
	logger       logger.Logger
}

func NewBot(cfg *config.Config, es *entrepreneur.Service, is *income.Service, log logger.Logger) (*Bot, error) {
	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	handlers := []handler{
		commands.NewStartHandler(),
		commands.NewHelpHandler(log),
		commands.NewRegisterHandler(es, log),
		commands.NewAddIncomeHandler(is),
		commands.NewGetIncomeHandler(is),
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
		logger:       log,
	}, nil
}

func (bot *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) {
	defer func() {
		if r := recover(); r != nil {
			bot.logger.Errorf(
				"panic recovered while handling update=%+v: %v\n%s",
				update,
				r,
				string(debug.Stack()),
			)
		}
	}()

	if update.Message == nil {
		return
	}
	h, ok := bot.cmdToHandler[update.Message.Command()]
	if !ok {
		// TODO: default handler
		return
	}
	_, err := h.Handle(ctx, bot.api, update.Message)
	if err != nil {
		bot.logger.Errorf(
			"failed to handle command=%q user_id=%d username=%q text=%q: %v",
			update.Message.Command(),
			update.Message.From.ID,
			update.Message.From.UserName,
			update.Message.Text,
			err,
		)
	}
}

func (bot *Bot) SendMessage(chatID uint, text string) error {
	msg := tgbotapi.NewMessage(int64(chatID), text)
	_, err := bot.api.Send(msg)
	if err != nil {
		bot.logger.Errorf("Failed to send message to chat %d: %v", chatID, err)
		return err
	}
	return nil
}

func (bot *Bot) Run(ctx context.Context) error {
	bot.logger.Info("bot started")
	for {
		select {
		case <-ctx.Done():
			bot.logger.Info("bot stopped")
			return nil
		case update := <-bot.updates:
			bot.handleUpdate(ctx, update)
		}
	}
}
