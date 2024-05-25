package commands

import (
	"fmt"
	"log"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"
	"paola-go-bot/telegram/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var BirthdayDelete = Command{
	Name:        "dimentica",
	Description: "Dimentica un compleanno",
	Handle:      handleBirthdayDelete,
}

func handleBirthdayDelete(message *tgbotapi.Message) status.CommandResponse {
	log.Printf("Delete birthdays - ask which one to delete")

	birthdays, ok := database.BirthdayFind(0, 0, message.Chat.ID)
	if !ok {
		log.Printf("Error deleting birthdays, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, ErrorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if len(birthdays) == 0 {
		log.Printf("Warning getting birthdays, no birthdays found yet")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Non ci sono compleanni ancora 🥲")
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	var rows [][]string = make([][]string, len(birthdays))

	for i, birthday := range birthdays {
		rows[i] = []string{birthday.Name}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, quale?")
	status.SetNext(message.Chat.ID, confirmDeletedBirthday)
	birthdaysKeyboard := utils.Keyboard(rows...)
	return status.CommandResponse{Reply: &reply, Keyboard: &birthdaysKeyboard}
}

func confirmDeletedBirthday(message *tgbotapi.Message, args ...string) status.CommandResponse {
	log.Printf("Delete birthdays - confirming deletion")

	name := message.Text

	response, err := DeleteBirthday(name, message.Chat.ID)
	if err != nil {
		log.Printf("Error deleting birthdays: %s", err)
		reply := tgbotapi.NewMessage(message.Chat.ID, ErrorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, response)
	status.ResetNext(message.Chat.ID)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}
}

func DeleteBirthday(name string, chatID int64) (string, error) {
	ok := database.BirthdayDeleteByName(name, chatID)
	if !ok {
		return "", fmt.Errorf("could not delete from database")
	}

	log.Printf("Birthday deleted: %s", name)
	return "Ok, me lo dimenticherò ✌️", nil
}
