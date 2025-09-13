package adapters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

type TelegramNotifier struct {
	bot *tgbotapi.BotAPI
}

func NewTelegramNotifier(bot *tgbotapi.BotAPI) *TelegramNotifier {
	return &TelegramNotifier{bot: bot}
}

func (t *TelegramNotifier) SendMessage(userID int64, message string) error {
	msg := tgbotapi.NewMessage(userID, message)

	go func() {
		defer func() {
			if r := recover(); r != nil {
			}
		}()

		_, err := t.bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}()
	return nil
}
