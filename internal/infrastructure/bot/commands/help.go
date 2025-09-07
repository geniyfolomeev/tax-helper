package commands

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpHandler struct{}

func NewHelpHandler() *HelpHandler {
	return &HelpHandler{}
}

func (h *HelpHandler) Command() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "help",
		Description: "Show this help message",
	}
}

func (h *HelpHandler) Handle(_ context.Context, api *tgbotapi.BotAPI, msg *tgbotapi.Message) (tgbotapi.Message, error) {
	commands, err := api.GetMyCommands()
	if err != nil {
		log.Println(err)
		return tgbotapi.Message{}, err
	}
	text := "📖 *Bot Commands Guide*\n\n"
	for _, cmd := range commands {
		text += "/" + cmd.Command + " – " + cmd.Description + "\n\n"
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "Markdown"
	return api.Send(reply)
}
