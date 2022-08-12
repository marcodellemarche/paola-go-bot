package telegram

import (
	"log"
	"paola-go-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Bot(token string, app_env string, db_secret string) {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		log.Panic(err)
	}

	if app_env != "prod" {
		bot.Debug = true
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	log.Printf("Commands: %+v", tgbotapi.NewSetMyCommands(commands...))

	update_config := tgbotapi.NewUpdate(0)
	update_config.Timeout = 60

	database.Initialize()

	statusMap, c := make(StatusMap), make(chan StatusUpdate)

	go manageStatus(statusMap, c)

	updates := bot.GetUpdatesChan(update_config)

	for update := range updates {
		if update.Message != nil {
			go handleUpdate(bot, update.Message, statusMap, c)
		}
	}
}

func manageStatus(
	status StatusMap,
	c <-chan StatusUpdate,
) {
	for update := range c {
		status[update.id] = UserStatusNew(update.next, update.args...)

		log.Printf("Status updated")
	}
}

func handleUpdate(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	status StatusMap,
	c chan<- StatusUpdate,
) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	userId := message.From.ID

	switch message.Command() {
	case commandRememberBirthday.Command:
		AskForName(bot, message, c)
	case commandGetBirthdays.Command:
		GetMyBirthdays(bot, message, c)
	default:
		{
			if userStatus, exists := status[userId]; exists {
				if userStatus.next != nil {
					userStatus.next(bot, message, c, userStatus.args...)
				} else {
					log.Printf("User ID %d has next callback set to nil", userId)
					defaultAnswer(bot, message, c)
				}
			} else {
				log.Printf("User ID %d still not present into status map", userId)
				defaultAnswer(bot, message, c)
			}
		}
	}
}

func defaultAnswer(bot *tgbotapi.BotAPI, message *tgbotapi.Message, c chan<- StatusUpdate,) {
	reply := tgbotapi.NewMessage(message.Chat.ID, randomInsult())
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = emptyKeyboard

	c <- StatusUpdateNew(message.From.ID, nil)

	bot.Send(reply)
}
