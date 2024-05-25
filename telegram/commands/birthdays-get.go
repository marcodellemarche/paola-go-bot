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

	list, err := ListBirthdays(message.Chat.ID)
	if err != nil {
		log.Printf("Error getting birthdays: %s", err)

		reply := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%s", err))
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, list)
	status.ResetNext(message.Chat.ID)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}
}

func ListBirthdays(chatID int64) (string, error) {
	birthdays, ok := database.BirthdayFind(0, 0, chatID)
	if !ok {
		return "", fmt.Errorf("could not fetch database with find")
	}

	birthdaysFromList, ok := database.BirthdayFindByList(0, 0, 0, chatID)
	if !ok {
		return "", fmt.Errorf("could not fetch database with find-by-list")
	}

	birthdays = append(birthdays, birthdaysFromList...)

	if len(birthdays) == 0 {
		return "Non ci sono compleanni ancora 🥲", nil
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

	message := "Tiè, ecco i compleanni:\n"

	for _, birthday := range birthdays {
		if birthday.UserName == "" {
			message += fmt.Sprintf("\n%s - %02d/%02d", birthday.Name, birthday.Day, birthday.Month)
		} else {
			message += fmt.Sprintf("\n[%s] %s - %02d/%02d", birthday.UserName, birthday.Name, birthday.Day, birthday.Month)
		}
	}

	return message, nil
}
