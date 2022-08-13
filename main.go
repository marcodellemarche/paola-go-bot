package main

import (
	"log"
	"strconv"

	"paola-go-bot/scripts"

	"github.com/alecthomas/kong"
	"github.com/joho/godotenv"
)

type Context struct {
	Debug bool
}

type ListenCmd struct {
	Super bool `help:"Enable SuperPaola mode."`
}

type ReminderCmd struct {
	Super bool `help:"Enable SuperPaola mode."`

	Days string `arg name:"days" help:"How many days from today for the reminder." default:"0" type:"number"`
}

var cli struct {
	Debug bool `help:"Enable debug mode."`

	Reminder ReminderCmd `cmd help:"Remember birthdays."`
	Listen ListenCmd   `cmd help:"Start the bot to listen for updates."`
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	ctx := kong.Parse(&cli)
	// Call the Run() method of the selected parsed command.
	err := ctx.Run(&Context{Debug: cli.Debug})
	ctx.FatalIfErrorf(err)
}

func (l *ListenCmd) Run(ctx *Context) error {
	log.Println("listen")
	scripts.Listener(ctx.Debug)

	return nil
}

func (r *ReminderCmd) Run(ctx *Context) error {
	log.Println("reminder", r.Days)
	days, _ := strconv.Atoi(r.Days)
	scripts.BirthdayReminder(days, ctx.Debug)

	return nil
}
