package scripts

import (
	"os"

	"paola-go-bot/database"
	"paola-go-bot/telegram"
)

func Listener(debug bool) {
	telegram_token := os.Getenv("TELEGRAM_TOKEN")
	database_url := os.Getenv("DATABASE_URL")

	telegram.Initialize(telegram_token, debug)

	database.Initialize(database_url, debug)

	telegram.ListenToUpdates()
}
