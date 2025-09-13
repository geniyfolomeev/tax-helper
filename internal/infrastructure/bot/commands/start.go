package commands

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartHandler struct{}

func NewStartHandler() *StartHandler {
	return &StartHandler{}
}

func (h *StartHandler) Command() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "start",
		Description: "Start interacting with the bot",
	}
}

func (h *StartHandler) Handle(_ context.Context, msg *tgbotapi.Message) tgbotapi.MessageConfig {
	return tgbotapi.NewMessage(msg.Chat.ID, "Привет! Я бот-помощник по налогам 🧾")
}
