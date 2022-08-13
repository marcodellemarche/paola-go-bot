package telegram

import (
	"log"

	"paola-go-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AskWhichToForget(
	message *tgbotapi.Message,
	c chan<- StatusUpdate,
) {
	log.Printf("Delete birthdays - ask which one to delete")

	birthdays, ok := database.BirthdayFindByUser(message.From.ID)
	if !ok {
		log.Printf("Error deleting birthdays, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, "So 'ncazzo io, ma qualcosa è andato storto 🥲")
		reply.ReplyToMessageID = message.MessageID
		reply.ReplyMarkup = emptyKeyboard
		c <- StatusUpdateNew(message.From.ID, nil)
		bot.Send(reply)
		return
	}

	var rows [][]string = make([][]string, len(birthdays))

	for i, birthday := range birthdays {
		rows[i] = []string{birthday.Name}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, quale?")
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = keyboard(rows...)

	c <- StatusUpdateNew(message.From.ID, confirmDeletedBirthday)

	bot.Send(reply)
}

func confirmDeletedBirthday(
	message *tgbotapi.Message,
	c chan<- StatusUpdate,
	args ...string,
) {
	log.Printf("Delete birthdays - confirming deletion")

	name := message.Text

	ok := database.BirthdayDeleteByName(name, message.From.ID)
	if !ok {
		log.Printf("Error deleting birthdays, could not update database")
		reply := tgbotapi.NewMessage(message.Chat.ID, "So 'ncazzo io, ma qualcosa è andato storto 🥲")
		reply.ReplyToMessageID = message.MessageID
		reply.ReplyMarkup = emptyKeyboard
		c <- StatusUpdateNew(message.From.ID, nil)
		bot.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, me lo dimenticherò ✌️")
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = emptyKeyboard

	c <- StatusUpdateNew(message.From.ID, nil)

	bot.Send(reply)
}
