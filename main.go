package main

import (
	"log"
	"os"
	"paola-go-bot/hello"
	"paola-go-bot/telegram"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	token := os.Getenv("TELEGRAM_TOKEN")
	app_env := os.Getenv("APP_ENV")
	db_secret := os.Getenv("FAUNADB_SECRET")

	hello.Hello()

	telegram.Bot(token, app_env, db_secret)
}
