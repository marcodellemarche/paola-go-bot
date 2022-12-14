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
	database_user := os.Getenv("POSTGRES_USER")
	database_db := os.Getenv("POSTGRES_DB")
	database_password := os.Getenv("POSTGRES_PASSWORD")

	telegram.Initialize(telegram_token, false)

	database.Initialize(database_user, database_db, database_password, false)

	users, _ := database.UserFindAll()

	for i, user := range users {
		log.Printf("User %d: %d - %s", i, user.Id, user.Name)

		name := telegram.GetNameFromUserId(user.Id)

		log.Printf("Found name: %s", name)
	}
}
