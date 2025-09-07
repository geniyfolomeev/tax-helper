package commands

import (
	"context"
	"strings"
	"tax-helper/internal/logger"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpHandler struct {
	logger logger.Logger
}

func NewHelpHandler(l logger.Logger) *HelpHandler {
	return &HelpHandler{logger: l}
}

func (h *HelpHandler) Command() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "help",
		Description: "Show this help message",
	}
}

func escapeMarkdownV2(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)
	return replacer.Replace(text)
}

func (h *HelpHandler) Handle(_ context.Context, api *tgbotapi.BotAPI, msg *tgbotapi.Message) (tgbotapi.Message, error) {
	commands, err := api.GetMyCommands()
	if err != nil {
		h.logger.Error(err)
		return tgbotapi.Message{}, err
	}
	text := "📖 *Bot Commands Guide*\n\n"
	for _, cmd := range commands {
		text += "/" + escapeMarkdownV2(cmd.Command) + " – " + escapeMarkdownV2(cmd.Description) + "\n\n"
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "MarkdownV2"
	return api.Send(reply)
}
