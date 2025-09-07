package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"tax-helper/internal/service/income"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const addIncomeExample = "/add_income DD.MM.YYYY AMOUNT CURRENCY"

type AddIncomeHandler struct {
	service *income.Service
}

func NewAddIncomeHandler(s *income.Service) *AddIncomeHandler {
	return &AddIncomeHandler{service: s}
}

func (h *AddIncomeHandler) Command() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "add_income",
		Description: "Add income",
	}
}

func (h *AddIncomeHandler) Handle(ctx context.Context, api *tgbotapi.BotAPI, msg *tgbotapi.Message) (tgbotapi.Message, error) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) != 3 {
		reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Wrong number of arguments, usage: %s", addIncomeExample))
		return api.Send(reply)
	}

	incomeDate, err := time.Parse("02.01.2006", args[0])
	if err != nil {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Wrong date format")
		return api.Send(reply)
	}

	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Wrong amount format")
		return api.Send(reply)
	}

	currency := args[2]

	err = h.service.AddIncome(ctx, uint(msg.Chat.ID), amount, currency, incomeDate)
	if err != nil {
		reply := tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		return api.Send(reply)
	}

	reply := tgbotapi.NewMessage(msg.Chat.ID, "Income added successfully")
	return api.Send(reply)
}
