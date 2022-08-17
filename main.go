package main

import (
	"log"
	"os"
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

type MigrationCmd struct {
	Write bool `help:"Write into database."`
	Token string `arg name:"token" help:"Telegram token." default:"" type:"string"`
}

type UsersCmd struct {}

var cli struct {
	Debug bool `help:"Enable debug mode."`

	Reminder  ReminderCmd  `cmd help:"Remember birthdays."`
	Listen    ListenCmd    `cmd help:"Start the bot to listen for updates."`
	Migration MigrationCmd `cmd help:"Start the DB migration."`
	Users UsersCmd `cmd help:"Fetch users info."`
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Println("Error loading .env file")
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
	
	os.Exit(0)

	return nil
}

func (m *MigrationCmd) Run(ctx *Context) error {
	log.Println("migration")
	scripts.Migration(m.Token, m.Write)

	return nil
}

func (l *UsersCmd) Run(ctx *Context) error {
	log.Println("users")
	scripts.Users()

	return nil
}
