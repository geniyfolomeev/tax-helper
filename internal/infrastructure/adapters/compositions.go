package adapters

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"tax-helper/internal/domain"
)

func CompositionsAdapters(bot *tgbotapi.BotAPI) map[string]domain.Notifier {
	return map[string]domain.Notifier{
		"tg":    NewTelegramNotifier(bot), // создаём инстанс
		"email": NewEmailNotifier(),
	}
}
