package telegram

import (
	"log"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AskWhichBirthdayToForget(
	message *tgbotapi.Message,
) {
	log.Printf("Delete birthdays - ask which one to delete")

	birthdays, ok := database.BirthdayFind(0, 0, message.Chat.ID)
	if !ok {
		log.Printf("Error deleting birthdays, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		SendMessage(reply, nil)
		return
	}

	if len(birthdays) == 0 {
		log.Printf("Warning getting birthdays, no birthdays found yet")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Non ci sono compleanni ancora 🥲")
		status.ResetNext(message.Chat.ID)
		SendMessage(reply, nil)
		return
	}

	var rows [][]string = make([][]string, len(birthdays))

	for i, birthday := range birthdays {
		rows[i] = []string{birthday.Name}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, quale?")

	status.SetNext(message.Chat.ID, confirmDeletedBirthday)

	birthdaysKeyboard := keyboard(rows...)

	SendMessage(reply, &birthdaysKeyboard)
}

func confirmDeletedBirthday(
	message *tgbotapi.Message,
	args ...string,
) {
	log.Printf("Delete birthdays - confirming deletion")

	name := message.Text

	ok := database.BirthdayDeleteByName(name, message.Chat.ID)
	if !ok {
		log.Printf("Error deleting birthdays, could not update database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		SendMessage(reply, nil)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, me lo dimenticherò ✌️")

	status.ResetNext(message.Chat.ID)

	SendMessage(reply, nil)
}
