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

type addIncomeArgs struct {
	incomeDate time.Time
	amount     float64
	currency   string
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

func (h *AddIncomeHandler) parseArgs(args []string) (*addIncomeArgs, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("wrong number of arguments, usage: %s", addIncomeExample)
	}

	incomeDate, err := time.Parse("02.01.2006", args[0])
	if err != nil {
		return nil, fmt.Errorf("wrong date format")
	}

	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, fmt.Errorf("wrong amount format")
	}

	currency := args[2]
	if len(currency) != 3 {
		return nil, fmt.Errorf("currency code must be 3 letters")
	}

	return &addIncomeArgs{
		incomeDate: incomeDate,
		amount:     amount,
		currency:   strings.ToUpper(currency),
	}, nil
}

func (h *AddIncomeHandler) Handle(ctx context.Context, msg *tgbotapi.Message) tgbotapi.MessageConfig {
	args, err := h.parseArgs(strings.Fields(msg.CommandArguments()))
	if err != nil {
		return tgbotapi.NewMessage(msg.Chat.ID, err.Error())
	}

	err = h.service.AddIncome(ctx, uint(msg.Chat.ID), args.amount, args.currency, args.incomeDate)
	if err != nil {
		return tgbotapi.NewMessage(msg.Chat.ID, err.Error())
	}
	return tgbotapi.NewMessage(msg.Chat.ID, "Income added successfully")
}
