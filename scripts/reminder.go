package scripts

import (
	"fmt"
	"log"
	"os"
	"time"

	"paola-go-bot/database"
	"paola-go-bot/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func BirthdayReminder(days int, debug bool) {
	token := os.Getenv("TELEGRAM_TOKEN")
	database_uri := os.Getenv("DATABASE_URI")

	telegram.Initialize(token, debug)

	database.Initialize(database_uri, debug)

	date := time.Now().AddDate(0, 0, days)

	birthdays, ok := database.BirthdayFindByDate(uint8(date.Day()), uint8(date.Month()))
	if !ok {
		log.Printf("Error getting birthdays, could not fetch database")
		return
	}

	printableDay := fmt.Sprintf("Tra %d giorni", days)
	if days == 0 {
		printableDay = "Oggi"
	} else if days == 1 {
		printableDay = "Domani"
	}

	for _, birthday := range birthdays {
		log.Printf("Notifying %d of %s's birthday", birthday.UserId, birthday.Name)
		message := tgbotapi.NewMessage(birthday.UserId, fmt.Sprintf("%s è il compleanno di %s!", printableDay, birthday.Name))
		telegram.SendMessage(message)
	}
}
