package adapters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tax-helper/internal/logger"
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
				logger.Error("Паника при отправке сообщения в чат %d: %v", userID, r)
			}
		}()

		_, err := t.bot.Send(msg)
		if err != nil {
			logger.Error("Ошибка отправки сообщения в чат %d: %v", userID, err)
		}
	}()
	return nil
}
