package telegram

import (
	"fmt"
	"log"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AskWhichListToSubscribe(
	message *tgbotapi.Message,
) {
	log.Printf("Subscribe list - ask which one to subscribe to")

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, condividi il contatto Telegram del tuo amico")

	status.SetNext(message.Chat.ID, confirmListSubscription)

	SendMessage(reply, nil)
}

func confirmListSubscription(
	message *tgbotapi.Message,
	args ...string,
) {
	log.Printf("Subscribe list - confirming subscription")

	if message.Contact == nil {
		log.Printf("Error subscribing new list, no contact provided")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma non hai condiviso nessun contatto!")
		status.ResetNext(message.Chat.ID)
		SendMessage(reply, nil)
		return
	}

	if message.Contact.UserID == 0 {
		log.Printf("Error subscribing new list, provided contact is not on Telegram")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il tuo amico non ha Telegram!")
		status.ResetNext(message.Chat.ID)
		SendMessage(reply, nil)
		return
	}

	var name string
	if message.Contact.LastName == "" {
		name = message.Contact.FirstName
	} else {
		name = fmt.Sprintf("%s %s", message.Contact.FirstName, message.Contact.LastName)
	}

	ok := database.ListInsert(message.Contact.UserID, message.Chat.ID, name)
	if !ok {
		log.Printf("Error subscribing new list, could not update database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		SendMessage(reply, nil)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, ora riceverai anche i suoi compleanni ✌️")

	status.ResetNext(message.Chat.ID)

	SendMessage(reply, nil)
}
