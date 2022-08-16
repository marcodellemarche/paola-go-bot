package telegram

import (
	"fmt"
	"log"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CheckIfNewUser(
	message *tgbotapi.Message,
	userId int64,
) bool {
	user, ok := database.UserFindById(message.Chat.ID)
	if !ok {
		log.Printf("Error finding user, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		SendMessage(reply, nil)
		return false
	}

	if user.Id == 0 {
		reply := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Te sei nuovo, inizia un po' con /%s", commandStart.Command))
		status.ResetNext(message.Chat.ID)
		SendMessage(reply, nil)
		return false
	}

	if _, exists := status.Get(userId); !exists {
		status.ResetNext(message.Chat.ID)
	}

	return true
}

func StartUser(
	message *tgbotapi.Message,
) {
	log.Printf("Start user")

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
		SendMessage(reply, nil)
		return
	}

	welcomeMessage := fmt.Sprintf("Ciao %s! Mi chiamo Paola Bartolbot e ricordo cose, tipo i compleanni.", message.From.FirstName)

	reply := tgbotapi.NewMessage(message.Chat.ID, welcomeMessage)

	status.ResetNext(message.Chat.ID)

	SendMessage(reply, nil)
}
