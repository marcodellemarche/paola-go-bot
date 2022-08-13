package telegram

import (
	"fmt"
	"log"
	"paola-go-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMyBirthdays(
	message *tgbotapi.Message,
	c chan<- StatusUpdate,
) {
	log.Printf("Get birthdays")

	birthdays, ok := database.BirthdayFindByUser(message.From.ID)
	if !ok {
		log.Printf("Error getting birthdays, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, "So 'ncazzo io, ma qualcosa è andato storto 🥲")
		reply.ReplyToMessageID = message.MessageID
		reply.ReplyMarkup = emptyKeyboard
		c <- StatusUpdateNew(message.From.ID, nil)
		bot.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Tiè, ecco i compleanni:\n")
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = emptyKeyboard

	for _, birthday := range birthdays {
		reply.Text += fmt.Sprintf("\n%s - %02d/%02d", birthday.Name, birthday.Day, birthday.Month)
	}

	c <- StatusUpdateNew(message.From.ID, nil)

	bot.Send(reply)
}
