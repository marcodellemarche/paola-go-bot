package commands

import (
	"paola-go-bot/telegram/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ErrorMessage = "So 'ncazzo io, ma qualcosa è andato storto 🥲"

type Command struct {
	Name        string
	Description string
	Handle      func(message *tgbotapi.Message) status.CommandResponse
}

var CommandsEnabled = []tgbotapi.BotCommand{
	{
		Command:     Start.Name,
		Description: Start.Description,
	},
	{
		Command:     BirthdaySet.Name,
		Description: BirthdaySet.Description,
	},
	{
		Command:     BirthdayDelete.Name,
		Description: BirthdayDelete.Description,
	},
	{
		Command:     BirthdaysGet.Name,
		Description: BirthdaysGet.Description,
	},
	// {
	// 	Command: SubscribeList.Name,
	// 	Description: SubscribeList.Description,
	// },
	{
		Command:     Stop.Name,
		Description: Stop.Description,
	},
}
