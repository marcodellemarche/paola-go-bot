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

func BirthdayReminder(date time.Time) {
	token := os.Getenv("TELEGRAM_TOKEN")
	app_env := os.Getenv("APP_ENV")
	db_secret := os.Getenv("FAUNADB_SECRET")

	telegram.Initialize(token, app_env != "prod", db_secret)

	database.Initialize(app_env != "prod")

	birthdays, ok := database.BirthdayFindByDate(uint8(date.Day()), uint8(date.Month()))
	if !ok {
		log.Printf("Error getting birthdays, could not fetch database")
		return
	}

	for _, birthday := range birthdays {
		log.Printf("Notifying %d of %s's birthday", birthday.UserId, birthday.Name)
		message := tgbotapi.NewMessage(birthday.UserId, fmt.Sprintf("Oggi è il compleanno di %s!", birthday.Name))
		telegram.SendMessage(message)
	}
}
