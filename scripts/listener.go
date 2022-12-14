package scripts

import (
	"os"

	"paola-go-bot/database"
	"paola-go-bot/telegram"
)

func Listener(debug bool) {
	telegram_token := os.Getenv("TELEGRAM_TOKEN")
	database_user := os.Getenv("POSTGRES_USER")
	database_db := os.Getenv("POSTGRES_DB")
	database_password := os.Getenv("POSTGRES_PASSWORD")

	telegram.Initialize(telegram_token, debug)

	database.Initialize(database_user, database_db, database_password, debug)

	telegram.ListenToUpdates()
}
