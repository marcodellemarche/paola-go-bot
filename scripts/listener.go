package scripts

import (
	"os"

	"paola-go-bot/chatgpt"
	"paola-go-bot/database"
	"paola-go-bot/telegram"
)

func Listener(debug bool) {
	telegram_token := os.Getenv("TELEGRAM_TOKEN")
	database_user := os.Getenv("POSTGRES_USER")
	database_db := os.Getenv("POSTGRES_DB")
	database_password := os.Getenv("POSTGRES_PASSWORD")
	openai_api_key := os.Getenv("OPENAI_API_KEY")
	openai_assistant_id := os.Getenv("OPENAI_ASSISTANT_ID")

	chatgptClient := chatgpt.New(openai_api_key, openai_assistant_id)

	telegramBot := telegram.New(telegram_token, chatgptClient, debug)

	database.Initialize(database_user, database_db, database_password, debug)

	telegramBot.ListenToUpdates()
}
