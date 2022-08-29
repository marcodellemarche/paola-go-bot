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
	telegram_token := os.Getenv("TELEGRAM_TOKEN")
	database_url := os.Getenv("DATABASE_URL")

	telegram.Initialize(telegram_token, debug)

	database.Initialize(database_url, debug)

	date := time.Now().AddDate(0, 0, days)

	printableDay := fmt.Sprintf("Tra %d giorni", days)
	if days == 0 {
		printableDay = "Oggi"
	} else if days == 1 {
		printableDay = "Domani"
	}

	birthdays, ok := database.BirthdayFind(uint8(date.Day()), uint8(date.Month()), 0)
	if !ok {
		log.Printf("Error getting birthdays, could not fetch database")
		return
	}

	for _, birthday := range birthdays {
		log.Printf("Notifying %d of %s's birthday", birthday.UserId, birthday.Name)
		message := tgbotapi.NewMessage(birthday.UserId, fmt.Sprintf("%s è il compleanno di %s!", printableDay, birthday.Name))
		telegram.SendMessage(&message, nil)
	}

	birthdays, ok = database.BirthdayFindByList(uint8(date.Day()), uint8(date.Month()), 0, 0)
	if !ok {
		log.Printf("Error getting birthdays from list, could not fetch database")
		return
	}

	for _, birthday := range birthdays {
		log.Printf("[%s] Notifying %d of %s's birthday", birthday.ListName, birthday.UserId, birthday.Name)
		message := tgbotapi.NewMessage(birthday.UserId, fmt.Sprintf("[%s] %s è il compleanno di %s!", birthday.ListName, printableDay, birthday.Name))
		telegram.SendMessage(&message, nil)
	}
}
