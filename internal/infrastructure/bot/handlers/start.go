package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleStart(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, "Привет! Я бот-помощник по налогам 🧾")
	_, err := api.Send(reply)
	if err != nil {
		// TODO: норм логгер
		log.Println(err)
	}
}
