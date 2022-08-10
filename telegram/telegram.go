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
	// for {
	// 	select {
	// 	case update := <-c:
	// 		{
	// 			if entry, ok := status[update.id]; ok {
	// 				log.Printf("Updating: %s", update.command)
	// 				// Then we modify the copy
	// 				entry.command = update.command

	// 				// Then we reassign map entry
	// 				status[update.id] = entry
	// 			} else {
	// 				status[update.id] = userStatusNew(update.command, 0, 0)
	// 			}

	// 			log.Printf("Command set, status: %+v", status)
	// 		}
	// 	case update := <-m:
	// 		{
	// 			if entry, ok := status[update.id]; ok {
	// 				log.Printf("Updating: %d", update.month)
	// 				// Then we modify the copy
	// 				entry.month = update.month

	// 				// Then we reassign map entry
	// 				status[update.id] = entry
	// 			} else {
	// 				status[update.id] = userStatusNew("", update.month, 0)
	// 			}

	// 			log.Printf("Month set, status: %+v", status)
	// 		}
	// 	}
	// }

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

	if userStatus, exists := status[userId]; exists {
		if userStatus.next != nil {
			userStatus.next(bot, message, status, c, userStatus.args...)
			return
		} else {
			log.Printf("User ID %d has next callback set to nil", userId)
		}
	} else {
		log.Printf("User ID %d still not present into status map", userId)
	}

	switch message.Text {
	// case "month":
	// 	reply.ReplyMarkup = monthKeyboard
	// case "day":
	// 	reply.ReplyMarkup = dayKeyboard
	// case "last":
	// 	{
	// 		if userStatus, exists := status[userId]; exists {
	// 			reply.Text = "Last command: " + userStatus.command
	// 		} else {
	// 			log.Printf("User ID %d still not present into status map", userId)
	// 		}
	// 	}
	// case "close":
	// 	reply.ReplyMarkup = emptyKeyboard
	case "ricorda":
		AskForName(bot, message, status, c)
	default:
		defaultAnswer(bot, message)
	}
}

func defaultAnswer(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(message.Chat.ID, randomInsult())
	reply.ReplyToMessageID = message.MessageID

	bot.Send(reply)
}
