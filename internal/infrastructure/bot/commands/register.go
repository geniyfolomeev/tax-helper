package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tax-helper/internal/domain"
	"tax-helper/internal/logger"
	"tax-helper/internal/service/entrepreneur"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const registerExample = "/register DD.MM.YYYY AMOUNT [LAST_SENT_DATE]"

type RegisterHandler struct {
	service *entrepreneur.Service
	logger  logger.Logger
}

type registerArgs struct {
	registeredAt time.Time
	amount       float64
	lastSentAt   time.Time
}

func NewRegisterHandler(s *entrepreneur.Service, l logger.Logger) *RegisterHandler {
	return &RegisterHandler{service: s, logger: l}
}

func (h *RegisterHandler) Command() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command: "register",
		Description: fmt.Sprintf(
			"Register as entrepreneur.\nUsage: %s\n"+
				"- DD.MM.YYYY = registration date\n"+
				"- AMOUNT = total income so far (number)\n"+
				"- LAST_SENT_DATE = optional, last declaration date (DD.MM.YYYY)",
			registerExample,
		),
	}
}

func (h *RegisterHandler) parseArgs(args []string) (*registerArgs, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("wrong number of arguments, usage: %s", registerExample)
	}

	registeredAt, err := time.Parse("02.01.2006", args[0])
	if err != nil {
		return nil, fmt.Errorf("wrong date format")
	}

	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, fmt.Errorf("wrong amount format")
	}

	var lastSentAt time.Time
	if len(args) == 3 {
		t, err := time.Parse("02.01.2006", args[2])
		if err != nil {
			return nil, fmt.Errorf("wrong date format for last declaration date")
		}
		lastSentAt = t
	}
	return &registerArgs{
		registeredAt: registeredAt,
		amount:       amount,
		lastSentAt:   lastSentAt,
	}, nil
}

func (h *RegisterHandler) Handle(ctx context.Context, msg *tgbotapi.Message) tgbotapi.MessageConfig {
	args, err := h.parseArgs(strings.Fields(msg.CommandArguments()))
	if err != nil {
		return tgbotapi.NewMessage(msg.Chat.ID, err.Error())
	}

	err = h.service.CreateEntrepreneur(ctx, uint(msg.Chat.ID), args.registeredAt, args.lastSentAt, args.amount)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			return tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		}
		if errors.Is(err, domain.ErrEntrepreneurAlreadyExists) {
			return tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		}
		h.logger.Error(err.Error())
		return tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Something went totally wrong: %s", err.Error()))
	}
	return tgbotapi.NewMessage(msg.Chat.ID, "Successfully registered")
}
