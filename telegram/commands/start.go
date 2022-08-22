package commands

import (
	"fmt"
	"log"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Start = Command{
	Name:        "start",
	Description: "Iniziamo da qui",
	Handle:      handleStart,
}

func handleStart(message *tgbotapi.Message) status.CommandResponse {
	log.Println("Start command")

	var name string
	if message.From.LastName == "" {
		name = message.From.FirstName
	} else {
		name = fmt.Sprintf("%s %s", message.From.FirstName, message.From.LastName)
	}

	ok := database.UserInsert(message.Chat.ID, name)
	if !ok {
		log.Printf("Error creating new user, could not update database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	welcomeMessage := fmt.Sprintf("Ciao %s! Mi chiamo Paola Bartolbot e ricordo cose, tipo i compleanni.", message.From.FirstName)

	reply := tgbotapi.NewMessage(message.Chat.ID, welcomeMessage)
	status.ResetNext(message.Chat.ID)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}
}
