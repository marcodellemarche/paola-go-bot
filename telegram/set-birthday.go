package telegram

import (
	"log"
	"strconv"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AskForName(
	message *tgbotapi.Message,
	args ...string,
) {
	log.Printf("Set birthday - asking for name")

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, di chi è il compleanno? (puoi pure condividere il contatto)")

	status.SetNext(message.From.ID, askForMonth)

	SendMessage(reply, nil)
}

func askForMonth(
	message *tgbotapi.Message,
	args ...string,
) {
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

	status.SetNext(message.From.ID, askForDay, name, contactId)

	SendMessage(reply, &monthKeyboard)
}

func askForDay(
	message *tgbotapi.Message,
	args ...string,
) {
	log.Printf("Set birthday - received month, asking for day")

	month := message.Text

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, che giorno?")

	status.SetNext(message.From.ID, confirmNewBirthday, month)

	SendMessage(reply, &dayKeyboard)
}

func confirmNewBirthday(
	message *tgbotapi.Message,
	args ...string,
) {
	log.Printf("Set birthday - received day, confirming birthday")

	name := args[0]
	contactId := args[1]
	month := args[2]
	day := message.Text

	if name == "" {
		log.Printf("Error confirming birthday, name is not valid: <empty-string>")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il nome non è valido!")
		status.ResetNext(message.From.ID)
		SendMessage(reply, nil)
		return
	}

	parsedMonth, err := strconv.ParseUint(month, 10, 8)
	if err != nil {
		log.Printf("Error confirming birthday, month is not valid: %s", month)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il mese non è valido!")
		status.ResetNext(message.From.ID)
		SendMessage(reply, nil)
		return
	}

	parsedDay, err := strconv.ParseUint(day, 10, 8)
	if err != nil {
		log.Printf("Error confirming birthday, day is not valid: %s", day)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il giorno non è valido!")
		status.ResetNext(message.From.ID)
		SendMessage(reply, nil)
		return
	}

	var parsedContactId int64
	if contactId != "" {
		parsedContactId, err = strconv.ParseInt(contactId, 10, 64)
		if err != nil {
			log.Printf("Error confirming birthday, contact ID is not valid: %s", contactId)
			reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il contatto non è valido!")
			status.ResetNext(message.From.ID)
			SendMessage(reply, nil)
			return
		}
	}

	ok := database.BirthdayInsert(name, parsedContactId, uint8(parsedDay), uint8(parsedMonth), message.From.ID)
	if !ok {
		log.Printf("Error confirming birthday, could not update database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.From.ID)
		SendMessage(reply, nil)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, me lo ricorderò ✌️")

	status.ResetNext(message.From.ID)

	SendMessage(reply, nil)
}
