package commands

import (
	"tax-helper/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type AddIncomeHandler struct {
	service *service.TaxService
}

func NewAddIncomeHandler(s *service.TaxService) *AddIncomeHandler {
	return &AddIncomeHandler{service: s}
}

func (h *AddIncomeHandler) Command() tgbotapi.BotCommand {
	return tgbotapi.BotCommand{
		Command:     "add_income",
		Description: "Add income",
	}
}

func (h *AddIncomeHandler) Handle(api *tgbotapi.BotAPI, msg *tgbotapi.Message) (tgbotapi.Message, error) {
	//_, err := h.service.EntrepreneurRepo.GetByID(uint(msg.Chat.ID))
	//if err != nil {
	//	if errors.Is(err, domain.ErrEntrepreneurNotFound) {
	//		reply := tgbotapi.NewMessage(msg.Chat.ID, "You need to /register first")
	//		return api.Send(reply)
	//	}
	//	logger.Error(err.Error())
	//	reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Something went totally wrong: %s", err.Error()))
	//	return api.Send(reply)
	//}
	return tgbotapi.Message{}, nil
}
