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
	text := "👋 Hello!\n\n" +
		"I'm your personal tax helper bot for individual entrepreneurs. Here's what I can do for you:\n\n" +
		"📌 Remind you to submit tax declarations on time\n" +
		"📌 Keep track of your income and submitted declarations\n" +
		"📌 Parse bank statements from TBC and BOG\n" +
		"📌 Automatically convert your income into GEL\n\n" +
		"Type /help anytime to see all commands.\n\n" +
		"Ready to simplify your life? Start with /register 🚀"
	return tgbotapi.NewMessage(msg.Chat.ID, text)
}
