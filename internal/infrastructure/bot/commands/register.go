package commands

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tax-helper/internal/domain"
	"tax-helper/internal/logger"
	"tax-helper/internal/service"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const usageExample = "/register DD.MM.YYYY AMOUNT"

type RegisterHandler struct {
	service *service.EntrepreneurService
}

func NewRegisterHandler(s *service.EntrepreneurService) *RegisterHandler {
	return &RegisterHandler{service: s}
}

func (h *RegisterHandler) Command() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "register",
		Description: fmt.Sprintf("Start registration, usage: %s", usageExample),
	}
}

func (h *RegisterHandler) Handle(api *tgbotapi.BotAPI, msg *tgbotapi.Message) (tgbotapi.Message, error) {
	_, err := h.service.EntrepreneurRepo.GetByID(uint(msg.Chat.ID))
	if err == nil || !errors.Is(err, domain.ErrEntrepreneurNotFound) {
		reply := tgbotapi.NewMessage(msg.Chat.ID, "You are already registered")
		return api.Send(reply)
	}

	args := strings.Fields(msg.CommandArguments())
	if len(args) < 2 || len(args) > 3 {
		reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Wrong number of arguments, usage: %s", usageExample))
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

	entrepreneur := domain.NewEntrepreneur(uint(msg.Chat.ID), "active", registeredAt, lastSentAt, amount)
	err = h.service.CreateEntrepreneur(entrepreneur)

	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			reply := tgbotapi.NewMessage(msg.Chat.ID, err.Error())
			return api.Send(reply)
		}
		logger.Error(err.Error())
		reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Something went totally wrong: %s", err.Error()))
		return api.Send(reply)
	}
	reply := tgbotapi.NewMessage(msg.Chat.ID, "Successfully registered")
	return api.Send(reply)
}
