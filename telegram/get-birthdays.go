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

	birthdays, ok := database.BirthdayFindByUser(message.From.ID)
	if !ok {
		log.Printf("Error getting birthdays, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.From.ID)
		SendMessage(reply, nil)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Tiè, ecco i compleanni:\n")

	for _, birthday := range birthdays {
		reply.Text += fmt.Sprintf("\n%s - %02d/%02d", birthday.Name, birthday.Day, birthday.Month)
	}

	status.ResetNext(message.From.ID)

	SendMessage(reply, nil)
}
