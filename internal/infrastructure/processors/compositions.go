package processors

import (
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/bot"
)

func CompositionsAdapters(bot *bot.Bot) map[string]domain.Notifier {
	return map[string]domain.Notifier{
		"send_declaration": NewSendDeclarationProcessor(bot), // создаём инстанс
		"add_income":       NewEmailNotifier(),
	}
}
