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

func (h *RegisterHandler) Handle(ctx context.Context, api *tgbotapi.BotAPI, msg *tgbotapi.Message) (tgbotapi.Message, error) {
	args := strings.Fields(msg.CommandArguments())
	if len(args) < 2 || len(args) > 3 {
		reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Wrong number of arguments, usage: %s", registerExample))
		return api.Send(reply)
	}

	registeredAt, err := time.Parse("02.01.2006", args[0])
	if err != nil {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Wrong date format")
		return api.Send(reply)
	}

	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Wrong amount format")
		return api.Send(reply)
	}

	var lastSentAt time.Time
	if len(args) == 3 {
		t, err := time.Parse("02.01.2006", args[2])
		if err != nil {
			reply := tgbotapi.NewMessage(msg.Chat.ID, "Wrong date format for last declaration date")
			return api.Send(reply)
		}
		lastSentAt = t
	}

	err = h.service.CreateEntrepreneur(ctx, uint(msg.Chat.ID), registeredAt, lastSentAt, amount)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			reply := tgbotapi.NewMessage(msg.Chat.ID, err.Error())
			return api.Send(reply)
		}
		if errors.Is(err, domain.ErrEntrepreneurAlreadyExists) {
			reply := tgbotapi.NewMessage(msg.Chat.ID, err.Error())
			return api.Send(reply)
		}
		h.logger.Error(err.Error())
		reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Something went totally wrong: %s", err.Error()))
		return api.Send(reply)
	}
	reply := tgbotapi.NewMessage(msg.Chat.ID, "Successfully registered")
	return api.Send(reply)
}
