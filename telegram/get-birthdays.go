package telegram

import (
	"fmt"
	"log"
	"paola-go-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMyBirthdays(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	c chan<- StatusUpdate,
) {
	log.Printf("Get birthdays")

	tgbotapi.NewSetMyCommands()

	reply := tgbotapi.NewMessage(message.Chat.ID, "Tiè, ecco i tuoi compleanni:")
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = emptyKeyboard

	birthdays, retrieved := database.BirthdayFindByUser(message.From.ID)
	if !retrieved {
		log.Printf("Error getting birthdays, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, "So 'ncazzo io, ma qualcosa è andato storto!")
		c <- StatusUpdateNew(message.From.ID, nil)
		bot.Send(reply)
		return
	}

	for _, birthday := range birthdays {
		reply.Text += fmt.Sprintf("\n%s - %02d/%02d", birthday.Name, birthday.Day, birthday.Month)
	}

	c <- StatusUpdateNew(message.From.ID, nil)

	bot.Send(reply)
}
