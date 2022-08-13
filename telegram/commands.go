package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var commandStart = tgbotapi.BotCommand{
	Command: "start",
	Description: "Iniziamo da qui",
}

var commandRememberBirthday = tgbotapi.BotCommand{
	Command: "ricorda",
	Description: "Ricorda un compleanno",
}

var commandForgetBirthday = tgbotapi.BotCommand{
	Command: "dimentica",
	Description: "Dimentica un compleanno",
}

var commandGetBirthdays = tgbotapi.BotCommand{
	Command: "compleanni",
	Description: "Lista dei compleanni da ricordare",
}

var commandStop = tgbotapi.BotCommand{
	Command: "stop",
	Description: "Interrompi il comando attuale",
}

var commands = []tgbotapi.BotCommand{
	commandStart,
	commandRememberBirthday,
	commandForgetBirthday,
	commandGetBirthdays,
	commandStop,
}

/*
start - Iniziamo da qui
ricorda - Ricorda un compleanno
dimentica - Dimentica un compleanno
compleanni - Lista dei compleanni da ricordare
stop - Interrompi il comando attuale
*/