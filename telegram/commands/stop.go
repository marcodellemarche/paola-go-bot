package commands

import (
	"log"

	"paola-go-bot/telegram/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Stop = Command{
	Name:        "stop",
	Description: "Interrompi il comando attuale",
	Handle:      handleStop,
}

func handleStop(message *tgbotapi.Message) status.CommandResponse {
	log.Println("Stop command")

	status.ResetNext(message.Chat.ID)

	return status.CommandResponse{Reply: nil, Keyboard: nil}
}
