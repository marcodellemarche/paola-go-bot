package telegram

import (
	"fmt"
	"log"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMyBirthdays(
	message *tgbotapi.Message,
) {
	log.Printf("Get birthdays")

	birthdays, ok := database.BirthdayFind(0, 0, message.Chat.ID)
	if !ok {
		log.Printf("Error getting birthdays, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		SendMessage(reply, nil)
		return
	}

	birthdaysFromList, ok := database.BirthdayFindByList(0, 0, 0, message.Chat.ID)
	if !ok {
		log.Printf("Error getting birthdays from list, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		SendMessage(reply, nil)
		return
	}

	if len(birthdays) == 0 && len(birthdaysFromList) == 0 {
		log.Printf("Warning getting birthdays, no birthdays found yet")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Non ci sono compleanni ancora 🥲")
		status.ResetNext(message.Chat.ID)
		SendMessage(reply, nil)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Tiè, ecco i compleanni:\n")

	for _, birthday := range birthdays {
		reply.Text += fmt.Sprintf("\n%s - %02d/%02d", birthday.Name, birthday.Day, birthday.Month)
	}

	for _, birthday := range birthdaysFromList {
		reply.Text += fmt.Sprintf("\n[%s] %s - %02d/%02d", birthday.UserName, birthday.Name, birthday.Day, birthday.Month)
	}

	status.ResetNext(message.Chat.ID)

	SendMessage(reply, nil)
}
