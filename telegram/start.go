package telegram

import (
	"fmt"
	"log"
	"paola-go-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CheckIfNewUser(
	message *tgbotapi.Message,
	status StatusMap,
	c chan<- StatusUpdate,
	userId int64,
) bool {
	user, ok := database.UserFindById(message.From.ID)
	if !ok {
		log.Printf("Error finding user, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, "So 'ncazzo io, ma qualcosa è andato storto 🥲")
		reply.ReplyToMessageID = message.MessageID
		reply.ReplyMarkup = emptyKeyboard
		c <- StatusUpdateNew(message.From.ID, nil)
		bot.Send(reply)
		return false
	}

	if user.Id == 0 {
		reply := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Te sei nuovo, inizia un po' con /%s", commandStart.Command))
		reply.ReplyToMessageID = message.MessageID
		reply.ReplyMarkup = emptyKeyboard
		c <- StatusUpdateNew(message.From.ID, nil)
		bot.Send(reply)
		return false
	}

	if _, exists := status[userId]; !exists {
		c <- StatusUpdateNew(message.From.ID, nil)
	}

	return true
}

func StartUser(
	message *tgbotapi.Message,
	c chan<- StatusUpdate,
) {
	log.Printf("Start user")

	ok := database.UserInsert(message.From.ID, message.From.FirstName)
	if !ok {
		log.Printf("Error creating new user, could not updte database")
		reply := tgbotapi.NewMessage(message.Chat.ID, "So 'ncazzo io, ma qualcosa è andato storto 🥲")
		reply.ReplyToMessageID = message.MessageID
		reply.ReplyMarkup = emptyKeyboard
		c <- StatusUpdateNew(message.From.ID, nil)
		bot.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Ciao %s!", message.From.FirstName))
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = emptyKeyboard

	c <- StatusUpdateNew(message.From.ID, nil)

	bot.Send(reply)
}
