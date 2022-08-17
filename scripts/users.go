package scripts

import (
	_ "embed"
	"log"
	"os"

	"paola-go-bot/database"
	"paola-go-bot/telegram"
)

func Users() {
	token := os.Getenv("TELEGRAM_TOKEN")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	database.Initialize(DATABASE_URL, false)

	users, _ := database.UserFindAll()

	for i, user := range users {
		log.Printf("User %d: %d - %s", i, user.Id, user.Name)

		telegram.Initialize(token, false)

		name := telegram.GetNameFromUserId(user.Id)

		log.Printf("Found name: %s", name)
	}
}
