package scripts

import (
	"os"

	"paola-go-bot/database"
	"paola-go-bot/telegram"
)

func Listener(debug bool) {
	token := os.Getenv("TELEGRAM_TOKEN")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	telegram.Initialize(token, debug)

	database.Initialize(DATABASE_URL, debug)

	telegram.ListenToUpdates()
}
