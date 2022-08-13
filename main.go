package main

import (
	"log"
	"time"

	"paola-go-bot/scripts"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

var CLI struct {
	Listen struct {
		SayHi bool `help:"Say hi."`
	} `cmd help:"Remove files."`

	Reminder struct {
		SayHi bool `help:"Say hi."`
	} `cmd help:"Remind birthdays."`
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx := kong.Parse(&CLI)
	switch ctx.Command() {
	case "listen":
		{
			log.Println("Running listener")
			scripts.Listener()
		}
	case "reminder":
		{
			log.Println("Running reminder")
			scripts.BirthdayReminder(time.Now())
		}
	default:
		panic(ctx.Command())
	}
}
