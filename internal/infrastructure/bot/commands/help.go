package commands

import (
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
		Description: "see help message",
	}
}

func (h *HelpHandler) Handle(api *tgbotapi.BotAPI, msg *tgbotapi.Message) (tgbotapi.Message, error) {
	commands, err := api.GetMyCommands()
	if err != nil {
		log.Println(err)
		return tgbotapi.Message{}, err
	}
	text := "📖 *Справка по командам*\n\n"
	for _, cmd := range commands {
		text += "/" + cmd.Command + " – " + cmd.Description + "\n"
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "Markdown"
	return api.Send(reply)
}
