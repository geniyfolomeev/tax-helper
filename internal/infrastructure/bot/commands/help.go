package commands

import (
	"context"
	"tax-helper/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type commandsGetter interface {
	GetMyCommands() ([]tgbotapi.BotCommand, error)
}

type HelpHandler struct {
	api    commandsGetter
	logger logger.Logger
}

func NewHelpHandler(l logger.Logger, api commandsGetter) *HelpHandler {
	return &HelpHandler{logger: l, api: api}
}

func (h *HelpHandler) Command() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "help",
		Description: "Show this help message",
	}
}

func (h *HelpHandler) Handle(_ context.Context, msg *tgbotapi.Message) tgbotapi.MessageConfig {
	commands, err := h.api.GetMyCommands()
	if err != nil {
		h.logger.Error(err)
		return tgbotapi.NewMessage(msg.Chat.ID, "something went wrong, try again later")
	}
	text := "📖 *Bot Commands Guide*\n\n"
	for _, cmd := range commands {
		text += "/" + cmd.Command + " – " + cmd.Description + "\n\n"
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, escapeMarkdownV2(text))
	reply.ParseMode = "MarkdownV2"
	return reply
}
