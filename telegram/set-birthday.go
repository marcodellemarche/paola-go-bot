package telegram

import (
	"log"
	"paola-go-bot/database"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AskForName(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	c chan<- StatusUpdate,
	args ...string,
) {
	log.Printf("Set birthday - asking for name")

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, che nome?")
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = emptyKeyboard

	c <- StatusUpdateNew(message.From.ID, askForMonth)

	bot.Send(reply)
}

func askForMonth(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	c chan<- StatusUpdate,
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
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = monthKeyboard

	c <- StatusUpdateNew(message.From.ID, askForDay, name, contactId)

	bot.Send(reply)
}

func askForDay(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	c chan<- StatusUpdate,
	args ...string,
) {
	log.Printf("Set birthday - received month, asking for day")

	month := message.Text

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, che giorno?")
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = dayKeyboard

	c <- StatusUpdateNew(message.From.ID, confirmBirthday, append(args, month)...)

	bot.Send(reply)
}

func confirmBirthday(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	c chan<- StatusUpdate,
	args ...string,
) {
	log.Printf("Set birthday - received day, confirming birthday")

	name := args[0]
	contactId := args[1]
	month := args[2]
	day := message.Text
	userName := message.From.FirstName

	if name == "" {
		log.Printf("Error confirming birthday, name is not valid: <empty-string>")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il nome non è valido!")
		c <- StatusUpdateNew(message.From.ID, nil)
		bot.Send(reply)
		return
	}

	parsedMonth, err := strconv.ParseUint(month, 10, 8)
	if err != nil {
		log.Printf("Error confirming birthday, month is not valid: %s", month)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il mese non è valido!")
		c <- StatusUpdateNew(message.From.ID, nil)
		bot.Send(reply)
		return
	}

	parsedDay, err := strconv.ParseUint(day, 10, 8)
	if err != nil {
		log.Printf("Error confirming birthday, day is not valid: %s", day)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il giorno non è valido!")
		c <- StatusUpdateNew(message.From.ID, nil)
		bot.Send(reply)
		return
	}

	var parsedContactId int64
	if contactId != "" {
		parsedContactId, err = strconv.ParseInt(contactId, 10, 64)
		if err != nil {
			log.Printf("Error confirming birthday, contact ID is not valid: %s", contactId)
			reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il contatto non è valido!")
			c <- StatusUpdateNew(message.From.ID, nil)
			bot.Send(reply)
			return
		}
	}

	inserted := database.BirthdayInsert(name, parsedContactId, uint8(parsedDay), uint8(parsedMonth), message.From.ID, userName)
	if !inserted {
		log.Printf("Error confirming birthday, could not update database")
		reply := tgbotapi.NewMessage(message.Chat.ID, "So 'ncazzo io, ma qualcosa è andato storto!")
		c <- StatusUpdateNew(message.From.ID, nil)
		bot.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Perfetto, me lo ricorderò!")
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = emptyKeyboard

	c <- StatusUpdateNew(message.From.ID, nil)

	bot.Send(reply)
}
