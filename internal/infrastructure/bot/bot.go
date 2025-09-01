package bot

import (
	"context"
	"tax-helper/config"
	"tax-helper/internal/infrastructure/bot/handlers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func NewBot(cfg *config.Config) (*Bot, error) {
	botApi, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	//botCfg := tgbotapi.SetMyCommandsConfig{
	//	Commands:     nil,
	//	Scope:        nil,
	//	LanguageCode: "",
	//}
	//_, err = botApi.Request(botCfg)
	//if err != nil {
	//	return nil, err
	//}
	return &Bot{
		api:     botApi,
		updates: botApi.GetUpdatesChan(u),
	}, nil
}

func (bot *Bot) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case update := <-bot.updates:
			if update.Message == nil {
				continue
			}
			switch update.Message.Command() {
			case "start":
				handlers.HandleStart(bot.api, update.Message)
			default:
				handlers.HandleDefault(bot.api, update.Message)
			}
		}
	}
}
