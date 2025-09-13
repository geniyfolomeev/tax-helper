package commands

import (
	"context"
	"fmt"
	"tax-helper/internal/service/income"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type GetIncomeHandler struct {
	service *income.Service
}

func NewGetIncomeHandler(s *income.Service) *GetIncomeHandler {
	return &GetIncomeHandler{service: s}
}

func (h *GetIncomeHandler) Command() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "get_income",
		Description: "Get income",
	}
}

func (h *GetIncomeHandler) Handle(ctx context.Context, api *tgbotapi.BotAPI, msg *tgbotapi.Message) (tgbotapi.Message, error) {
	actualIncome, err := h.service.GetActualIncome(ctx, uint(msg.Chat.ID))
	if err != nil {
		reply := tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		return api.Send(reply)
	}

	text := "*Your income*\n\n"

	if len(actualIncome.CurMonth) > 0 {
		text += "*Current month:*\n"
		for _, i := range actualIncome.CurMonth {
			amount, _ := i.Amount.Float64()
			srcAmount, _ := i.SourceAmount.Float64()
			text += fmt.Sprintf("%s — %.2f GEL; (%.2f %s)\n", i.Date.Format("02 Jan 2006"), amount, srcAmount, i.SourceCurrency)
		}
		text += "\n"
	} else {
		text += "No income for current month.\n\n"
	}

	if len(actualIncome.PrevMonth) > 0 {
		text += "*Previous month:*\n"
		for _, i := range actualIncome.PrevMonth {
			amount, _ := i.Amount.Float64()
			srcAmount, _ := i.SourceAmount.Float64()
			text += fmt.Sprintf("%s — %.2f GEL; (%.2f %s)\n", i.Date.Format("02 Jan 2006"), amount, srcAmount, i.SourceCurrency)
		}
	} else {
		text += "No income for previous month."
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, escapeMarkdownV2(text))
	reply.ParseMode = "MarkdownV2"
	return api.Send(reply)
}
