package telegram

import (
	"log"
	"paola-go-bot/database"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var monthKeyboard = keyboard(
	[]string{"1", "2", "3"},
	[]string{"4", "5", "6"},
	[]string{"7", "8", "9"},
	[]string{"10", "11", "12"},
)

var dayKeyboard = keyboard(
	[]string{"1", "2", "3", "4", "5", "6", "7"},
	[]string{"8", "9", "10", "11", "12", "13"},
	[]string{"14", "15", "16", "17", "18", "19"},
	[]string{"20", "21", "22", "23", "24", "25"},
	[]string{"26", "27", "28", "29", "30", "31"},
)

var emptyKeyboard = tgbotapi.NewRemoveKeyboard(true)

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

	database.Initialize(db_secret)

	database.Sql()

	statusMap, statusChan := make(StatusMap), make(chan StatusUpdate)

	go manageStatus(statusMap, statusChan)

	updates := bot.GetUpdatesChan(update_config)

	for update := range updates {
		if update.Message != nil {
			go handleUpdate(bot, update.Message, statusMap, statusChan)
		}
	}
}

func manageStatus(status StatusMap, c <-chan StatusUpdate) {
	for update := range c {
		status[update.id] = userStatusNew(update.text)

		log.Printf("Status: %+v", status)
	}

}

func handleUpdate(bot *tgbotapi.BotAPI, message *tgbotapi.Message, status StatusMap, c chan<- StatusUpdate) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	c <- statusUpdateNew(message.From.ID, message.Text)

	msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	msg.ReplyToMessageID = message.MessageID

	switch message.Text {
	case "month":
		msg.ReplyMarkup = monthKeyboard
	case "day":
		msg.ReplyMarkup = dayKeyboard
	case "last":
		{
			if userStatus, exists := status[message.From.ID]; exists {
				msg.Text = "Last: " + userStatus.text
			} else {
				log.Printf("User ID %d still not present into status map", message.From.ID)
			}
		}
	case "close":
		msg.ReplyMarkup = emptyKeyboard
	}

	bot.Send(msg)
}
