package scripts

import (
	"os"

	"paola-go-bot/database"
	"paola-go-bot/telegram"
)

func Listener(debug bool) {
	token := os.Getenv("TELEGRAM_TOKEN")
	database_uri := os.Getenv("DATABASE_URI")

	telegram.Initialize(token, debug)

	database.Initialize(database_uri, debug)

	telegram.ListenToUpdates()
}
