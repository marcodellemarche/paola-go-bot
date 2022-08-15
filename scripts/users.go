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
	database_uri := os.Getenv("DATABASE_URI")

	database.Initialize(database_uri, false)

	users, _ := database.UserFindAll()

	for i, user := range users {
		log.Printf("User %d: %d - %s", i, user.Id, user.Name)

		telegram.Initialize(token, false)

		name := telegram.GetChat(user.Id)

		log.Printf("Found name: %s", name)
	}
}
