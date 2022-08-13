package scripts

import (
	"os"

	"paola-go-bot/database"
	"paola-go-bot/telegram"
)

func Listener() {
	token := os.Getenv("TELEGRAM_TOKEN")
	app_env := os.Getenv("APP_ENV")
	db_secret := os.Getenv("FAUNADB_SECRET")

	telegram.Initialize(token, app_env != "prod", db_secret)

	database.Initialize(app_env != "prod")

	telegram.ListenToUpdates()
}
