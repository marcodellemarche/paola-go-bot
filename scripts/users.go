package scripts

import (
	_ "embed"
	"log"
	"os"

	"paola-go-bot/database"
	"paola-go-bot/telegram"
)

func Users() {
	telegram_token := os.Getenv("TELEGRAM_TOKEN")
	database_url := os.Getenv("DATABASE_URL")

	telegram.Initialize(telegram_token, false)

	database.Initialize(database_url, false)

	users, _ := database.UserFindAll()

	for i, user := range users {
		log.Printf("User %d: %d - %s", i, user.Id, user.Name)

		name := telegram.GetNameFromUserId(user.Id)

		log.Printf("Found name: %s", name)
	}
}
