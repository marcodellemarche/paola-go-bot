package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var commandRememberBirthday = tgbotapi.BotCommand{
	Command: "ricorda",
	Description: "Ricorda un compleanno",
}

var commandGetBirthdays = tgbotapi.BotCommand{
	Command: "compleanni",
	Description: "Vedi i compleanni salvati",
}

var commands = []tgbotapi.BotCommand{
	commandRememberBirthday,
	commandGetBirthdays,
} 
