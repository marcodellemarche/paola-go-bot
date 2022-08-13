package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var bot tgbotapi.BotAPI

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
	log.Printf("Listening updates to commands: %+v", tgbotapi.NewSetMyCommands(commands...))

	update_config := tgbotapi.NewUpdate(0)
	update_config.Timeout = 60

	statusMap, c := make(StatusMap), make(chan StatusUpdate)

	go manageStatus(statusMap, c)

	updates := bot.GetUpdatesChan(update_config)

	for update := range updates {
		if update.Message != nil {
			go handleUpdate(update.Message, statusMap, c)
		}
	}
}

func handleUpdate(
	message *tgbotapi.Message,
	status StatusMap,
	c chan<- StatusUpdate,
) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	userId := message.From.ID

	if message.Command() == commandStart.Command {
		StartUser(message, c)
		return
	}

	ok := CheckIfNewUser(message, status, c, userId)
	if !ok {
		return
	}

	switch message.Command() {
	case commandRememberBirthday.Command:
		AskForName(message, c)
	case commandGetBirthdays.Command:
		GetMyBirthdays(message, c)
	case commandForgetBirthday.Command:
		AskWhichToForget(message, c)
	default:
		{
			if userStatus, exists := status[userId]; exists {
				if userStatus.next != nil {
					userStatus.next(message, c, userStatus.args...)
				} else {
					log.Printf("User ID %d has next callback set to nil", userId)
					defaultAnswer(message, c)
				}
			} else {
				log.Printf("User ID %d still not present into status map", userId)
				defaultAnswer(message, c)
			}
		}
	}
}

func defaultAnswer(
	message *tgbotapi.Message,
	c chan<- StatusUpdate,
) {
	reply := tgbotapi.NewMessage(message.Chat.ID, randomInsult())
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = emptyKeyboard

	c <- StatusUpdateNew(message.From.ID, nil)

	bot.Send(reply)
}

func SendMessage(message tgbotapi.MessageConfig) {
	bot.Send(message)
}
