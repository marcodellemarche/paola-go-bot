package commands

import (
	"fmt"
	"log"
	"strconv"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"
	"paola-go-bot/telegram/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var BirthdaySet = Command{
	Name:        "ricorda",
	Description: "Ricorda un compleanno",
	Handle:      handleBirthdaySet,
}

func handleBirthdaySet(message *tgbotapi.Message) status.CommandResponse {
	log.Printf("Set birthday - asking for name")

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, di chi è il compleanno? (puoi anche condividere il contatto Telegram)")
	status.SetNext(message.Chat.ID, askForBirthdayMonth)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}
}

func askForBirthdayMonth(message *tgbotapi.Message, args ...string) status.CommandResponse {
	log.Printf("Set birthday - received name, asking for month")

	var name string
	var contactId string

	if message.Text != "" {
		name = message.Text
	} else if message.Contact != nil {
		name = message.Contact.FirstName

		if message.Contact.UserID != 0 {
			contactId = strconv.FormatInt(message.Contact.UserID, 10)
		}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, che mese?")
	status.SetNext(message.Chat.ID, askForBirthdayDay, name, contactId)
	return status.CommandResponse{Reply: &reply, Keyboard: &utils.MonthKeyboard}
}

func askForBirthdayDay(message *tgbotapi.Message, args ...string) status.CommandResponse {
	log.Printf("Set birthday - received month, asking for day")

	month := message.Text

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, che giorno?")
	status.SetNext(message.Chat.ID, confirmNewBirthday, month)
	return status.CommandResponse{Reply: &reply, Keyboard: &utils.DayKeyboard}
}

func confirmNewBirthday(message *tgbotapi.Message, args ...string) status.CommandResponse {
	log.Printf("Set birthday - received day, confirming birthday")

	name := args[0]
	contactId := args[1]
	month := args[2]
	day := message.Text

	if name == "" {
		log.Printf("Error confirming birthday, name is not valid: <empty-string>")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il nome non è valido!")
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	parsedMonth, err := strconv.ParseUint(month, 10, 8)
	if err != nil {
		log.Printf("Error confirming birthday, month is not valid: %s", month)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il mese non è valido!")
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	parsedDay, err := strconv.ParseUint(day, 10, 8)
	if err != nil {
		log.Printf("Error confirming birthday, day is not valid: %s", day)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il giorno non è valido!")
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	var parsedContactId int64
	if contactId != "" {
		parsedContactId, err = strconv.ParseInt(contactId, 10, 64)
		if err != nil {
			log.Printf("Error confirming birthday, contact ID is not valid: %s", contactId)
			reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il contatto non è valido!")
			status.ResetNext(message.Chat.ID)
			return status.CommandResponse{Reply: &reply, Keyboard: nil}
		}
	}

	ok := database.BirthdayInsert(name, parsedContactId, uint8(parsedDay), uint8(parsedMonth), message.Chat.ID)
	if !ok {
		log.Printf("Error confirming birthday, could not update database")
		reply := tgbotapi.NewMessage(message.Chat.ID, ErrorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	response, err := SetBirthday(name, parsedContactId, uint8(parsedDay), uint8(parsedMonth), message.Chat.ID)
	if err != nil {
		log.Printf("Error confirming birthday: %s", err)
		reply := tgbotapi.NewMessage(message.Chat.ID, ErrorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, response)
	status.ResetNext(message.Chat.ID)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}
}

func SetBirthday(name string, contactId int64, day uint8, month uint8, chatID int64) (string, error) {
	ok := database.BirthdayInsert(name, 0, day, month, chatID)
	if !ok {
		return "", fmt.Errorf("could not update database")
	}

	log.Printf("Set birthday for %s to %d-%d", name, month, day)
	return "Ok, me lo ricorderò ✌️", nil
}
