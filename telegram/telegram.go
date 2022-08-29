package telegram

import (
	"fmt"
	"log"

	"paola-go-bot/telegram/commands"
	"paola-go-bot/telegram/status"
	"paola-go-bot/telegram/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot tgbotapi.BotAPI

func Initialize(token string, debug bool) {
	if token == "" {
		log.Fatal("Missing Telegram token")
	}

	pBot, err := tgbotapi.NewBotAPI(token)
	bot = *pBot

	if err != nil {
		log.Panic(err)
	}

	if debug {
		bot.Debug = true
	}

	log.Printf("Bot authorized on account %s", bot.Self.UserName)
}

func ListenToUpdates() {
	bot.Request(tgbotapi.NewSetMyCommands(commands.CommandsEnabled...))

	log.Printf("Listening updates to commands: %+v", commands.CommandsEnabled)

	update_config := tgbotapi.NewUpdate(0)
	update_config.Timeout = 60

	go status.Manage()

	updates := bot.GetUpdatesChan(update_config)

	for update := range updates {
		if update.Message != nil {
			go handleUpdate(update.Message)
		}
	}
}

func handleUpdate(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	userId := message.Chat.ID

	var response status.CommandResponse

	if message.Command() == commands.Start.Name {
		response = commands.Start.Handle(message)
		SendMessage(response.Reply, response.Keyboard)
		return
	}

	if response = commands.CheckIfNewUser(message, userId); response.Reply != nil {
		SendMessage(response.Reply, response.Keyboard)
		return
	}

	switch message.Command() {
	case commands.BirthdaySet.Name:
		response = commands.BirthdaySet.Handle(message)
	case commands.BirthdaysGet.Name:
		response = commands.BirthdaysGet.Handle(message)
	case commands.BirthdayDelete.Name:
		response = commands.BirthdayDelete.Handle(message)
	case commands.ListSet.Name:
		response = commands.ListSet.Handle(message)
	case commands.Stop.Name:
		response = commands.Stop.Handle(message)
	case commands.WishlistSet.Name:
		response = commands.WishlistSet.Handle(message)
	case commands.WishlistsBirthdaysGet.Name:
		response = commands.WishlistsBirthdaysGet.Handle(message)
	default:
		response = commands.CheckNextActionOrDefault(message, userId)
	}

	SendMessage(response.Reply, response.Keyboard)
}

func SendMessage(message *tgbotapi.MessageConfig, keyboard *tgbotapi.ReplyKeyboardMarkup) {
	if message == nil {
		return
	}

	if keyboard == nil {
		message.ReplyMarkup = utils.EmptyKeyboard
	} else {
		message.ReplyMarkup = keyboard
	}

	bot.Send(message)
}

func GetNameFromUserId(userId int64) string {
	chat, err := bot.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: userId}})
	if err != nil {
		log.Printf("Error getting chat for user %d: %s", userId, err.Error())
		return ""
	}

	if chat.LastName == "" {
		return chat.FirstName
	}

	return fmt.Sprintf("%s %s", chat.FirstName, chat.LastName)
}
