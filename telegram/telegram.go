package telegram

import (
	"fmt"
	"log"

	"paola-go-bot/telegram/status"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot tgbotapi.BotAPI

var errorMessage = "So 'ncazzo io, ma qualcosa è andato storto 🥲"

func Initialize(token string, debug bool) {
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
	bot.Request(tgbotapi.NewSetMyCommands(commands...))

	log.Printf("Listening updates to commands: %+v", commands)

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

func handleUpdate(
	message *tgbotapi.Message,
) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	userId := message.Chat.ID

	if message.Command() == commandStart.Command {
		StartUser(message)
		return
	}

	ok := CheckIfNewUser(message, userId)
	if !ok {
		return
	}

	switch message.Command() {
	case commandRememberBirthday.Command:
		AskForBirthdayName(message)
	case commandGetBirthdays.Command:
		GetMyBirthdays(message)
	case commandForgetBirthday.Command:
		AskWhichBirthdayToForget(message)
	case commandSubscribeList.Command:
		AskWhichListToSubscribe(message)
	case commandStop.Command:
		Stop(message)
	default:
		CheckNextActionOrDefault(message, userId)
	}
}

func SendMessage(message tgbotapi.MessageConfig, keyboard *tgbotapi.ReplyKeyboardMarkup) {
	if keyboard == nil {
		message.ReplyMarkup = emptyKeyboard
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

func Stop(
	message *tgbotapi.Message,
) {
	status.ResetNext(message.Chat.ID)
}

func CheckNextActionOrDefault(
	message *tgbotapi.Message,
	userId int64,
) {
	if userStatus, exists := status.Get(userId); exists {
		if userStatus.Next != nil {
			userStatus.Next(message, userStatus.Args...)
			return
		} else {
			log.Printf("User ID %d has next callback set to nil", userId)
		}
	} else {
		log.Printf("User ID %d still not present into status map", userId)
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, randomInsult())

	status.ResetNext(message.Chat.ID)

	SendMessage(reply, nil)
}
