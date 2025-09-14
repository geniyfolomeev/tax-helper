package processors

import (
	"context"
	"fmt"
	"tax-helper/internal/domain"
	"tax-helper/internal/infrastructure/bot"
)

type SendDeclarationProcessor struct {
	bot *bot.Bot
}

func NewSendDeclarationProcessor(bot *bot.Bot) *SendDeclarationProcessor {
	return &SendDeclarationProcessor{bot: bot}
}

func (t *SendDeclarationProcessor) Process(ctx context.Context, task domain.Task) error {
	if task.TelegramID == 0 {
		return fmt.Errorf("у задачи %v не задан TelegramID", task)
	}
	if task.Status != "ready" {
		return fmt.Errorf("декларация для %d не готова (status=%s)", task.TelegramID, task.Status)
	}

	// бизнес-логика
	text := fmt.Sprintf("Уважаемый пользователь %d, декларация готова для отправки ✅", task.TelegramID)

	return t.bot.SendMessage(task.TelegramID, text)
}
