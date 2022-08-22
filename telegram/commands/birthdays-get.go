package commands

import (
	"fmt"
	"log"
	"sort"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var BirthdaysGet = Command{
	Name:        "compleanni",
	Description: "Lista dei compleanni da ricordare",
	Handle:      handleBirthdaysGet,
}

func handleBirthdaysGet(message *tgbotapi.Message) status.CommandResponse {
	log.Printf("Get birthdays")

	birthdays, ok := database.BirthdayFind(0, 0, message.Chat.ID)
	if !ok {
		log.Printf("Error getting birthdays, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	birthdaysFromList, ok := database.BirthdayFindByList(0, 0, 0, message.Chat.ID)
	if !ok {
		log.Printf("Error getting birthdays from list, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	birthdays = append(birthdays, birthdaysFromList...)

	if len(birthdays) == 0 {
		log.Printf("Warning getting birthdays, no birthdays found yet")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Non ci sono compleanni ancora 🥲")
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	sort.Slice(birthdays, func(i, j int) bool {
		return birthdays[i].Before(&birthdays[j])
	})

	for i, birthday := range birthdays {
		if !birthday.Passed() {
			birthdays = append(birthdays[i:], birthdays[:i]...)
			break
		}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Tiè, ecco i compleanni:\n")

	for _, birthday := range birthdays {
		if birthday.UserName == "" {
			reply.Text += fmt.Sprintf("\n%s - %02d/%02d", birthday.Name, birthday.Day, birthday.Month)
		} else {
			reply.Text += fmt.Sprintf("\n[%s] %s - %02d/%02d", birthday.UserName, birthday.Name, birthday.Day, birthday.Month)
		}
	}

	status.ResetNext(message.Chat.ID)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}
}
