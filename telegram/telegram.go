package telegram

import (
	"fmt"
	"log"

	"paola-go-bot/chatgpt"
	"paola-go-bot/database"
	"paola-go-bot/telegram/commands"
	"paola-go-bot/telegram/status"
	"paola-go-bot/telegram/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	bot           tgbotapi.BotAPI
	chatgptClient *chatgpt.ChatGPT
}

func New(token string, chatgptClient *chatgpt.ChatGPT, debug bool) *Telegram {
	if token == "" {
		log.Fatal("Missing Telegram token")
	}

	var bot tgbotapi.BotAPI
	pBot, err := tgbotapi.NewBotAPI(token)
	bot = *pBot

	if err != nil {
		log.Panic(err)
	}

	if debug {
		bot.Debug = true
	}

	log.Printf("Bot authorized on account %s", bot.Self.UserName)

	return &Telegram{
		bot:           bot,
		chatgptClient: chatgptClient,
	}
}

func (t *Telegram) ListenToUpdates() {
	t.bot.Request(tgbotapi.NewSetMyCommands(commands.CommandsEnabled...))

	log.Printf("Listening updates to commands: %+v", commands.CommandsEnabled)

	update_config := tgbotapi.NewUpdate(0)
	update_config.Timeout = 60

	go status.Manage()

	updates := t.bot.GetUpdatesChan(update_config)

	for update := range updates {
		if update.Message != nil {
			go t.handleUpdate(update.Message)
		}
	}
}

func (t *Telegram) handleUpdate(message *tgbotapi.Message) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	userId := message.Chat.ID

	var response status.CommandResponse

	if message.Command() == commands.Start.Name {
		response = commands.Start.Handle(message)
		t.SendMessage(response.Reply, response.Keyboard)
		return
	}

	if response = t.checkIfNewUser(message, userId); response.Reply != nil {
		t.SendMessage(response.Reply, response.Keyboard)
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
	default:
		response = t.checkNextActionOrDefault(message, userId)
	}

	t.SendMessage(response.Reply, response.Keyboard)
}

func (t *Telegram) SendMessage(message *tgbotapi.MessageConfig, keyboard *tgbotapi.ReplyKeyboardMarkup) {
	if message == nil {
		return
	}

	if keyboard == nil {
		message.ReplyMarkup = utils.EmptyKeyboard
	} else {
		message.ReplyMarkup = keyboard
	}

	t.bot.Send(message)
}

func (t *Telegram) GetNameFromUserId(userId int64) string {
	chat, err := t.bot.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: userId}})
	if err != nil {
		log.Printf("Error getting chat for user %d: %s", userId, err.Error())
		return ""
	}

	if chat.LastName == "" {
		return chat.FirstName
	}

	return fmt.Sprintf("%s %s", chat.FirstName, chat.LastName)
}

func (t *Telegram) checkIfNewUser(message *tgbotapi.Message, userId int64) status.CommandResponse {
	user, ok := database.UserFindById(message.Chat.ID)
	if !ok {
		log.Printf("Error finding user, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, commands.ErrorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if user.Id == 0 {
		reply := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Te sei nuovo, inizia un po' con /%s", commands.Start.Name))
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if _, exists := status.Get(userId); !exists {
		status.ResetNext(message.Chat.ID)
	}

	return status.CommandResponse{Reply: nil, Keyboard: nil}
}

func (t *Telegram) checkNextActionOrDefault(message *tgbotapi.Message, userId int64) status.CommandResponse {
	if userStatus, exists := status.Get(userId); exists {
		if userStatus.Next != nil {
			return userStatus.Next(message, userStatus.Args...)
		} else {
			log.Printf("User ID %d has next callback set to nil", userId)
		}
	} else {
		log.Printf("User ID %d still not present into status map", userId)
	}

	status.ResetNext(message.Chat.ID)

	if t.chatgptClient == nil {
		log.Printf("ChatGPT client is not setup, cannot use AI assistant")

		reply := tgbotapi.NewMessage(message.Chat.ID, utils.RandomInsult())
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if message.Voice != nil {
		return t.aiHandleVoice(message)
	} else {
		return t.aiHandleText(message)
	}
}
