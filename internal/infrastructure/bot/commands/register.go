//go:generate mockgen -source=$GOFILE -destination=mocks/mock_register.go -package=mocks
package commands

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"tax-helper/internal/domain"
	"tax-helper/internal/logger"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type entrepreneurCreator interface {
	CreateEntrepreneur(ctx context.Context, tgID int64, regAt, lastAt time.Time, yta float64) error
}

type RegisterHandler struct {
	service entrepreneurCreator
	logger  logger.Logger
}

type registerArgs struct {
	registeredAt time.Time
	amount       float64
	lastSentAt   time.Time
}

const registerUsage = "*Usage*:\n" +
	"/register  *date*  *amount*  last_declaration_date\n\n" +
	"*Arguments*:\n" +
	"  *date* — The date you registered as an entrepreneur (DD.MM.YYYY)\n" +
	"  *amount* — Total income so far this year in GEL\n" +
	"  *last_declaration_date* — Optional. The date of your last submitted declaration (DD.MM.YYYY or empty)\n\n" +
	"*Examples*:\n" +
	"  /register 31.01.2025 15000 31.08.2025\n" +
	"  /register 01.01.2025 0"

func NewRegisterHandler(s entrepreneurCreator, l logger.Logger) *RegisterHandler {
	return &RegisterHandler{service: s, logger: l}
}

func (h *RegisterHandler) Command() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "register",
		Description: "Register as entrepreneur",
	}
}

func (h *RegisterHandler) parseArgs(args []string) (*registerArgs, error) {
	if len(args) < 2 || len(args) > 3 {
		return nil, fmt.Errorf("❌ Wrong number of arguments")
	}

	registeredAt, err := time.Parse("02.01.2006", args[0])
	if err != nil {
		return nil, fmt.Errorf("❌ Wrong date format")
	}

	amount, err := strconv.ParseFloat(args[1], 64)
	if err != nil {
		return nil, fmt.Errorf("❌ Wrong amount format")
	}

	var lastSentAt time.Time
	if len(args) == 3 {
		t, err := time.Parse("02.01.2006", args[2])
		if err != nil {
			return nil, fmt.Errorf("❌ Wrong date format for last declaration date")
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
		reply := tgbotapi.NewMessage(msg.Chat.ID, escapeMarkdownV2(fmt.Sprintf("%s\n\n%s", err.Error(), registerUsage)))
		reply.ParseMode = "MarkdownV2"
		return reply
	}

	err = h.service.CreateEntrepreneur(ctx, msg.Chat.ID, args.registeredAt, args.lastSentAt, args.amount)
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
