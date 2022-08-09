package telegram

import (
	"log"
	"paola-go-bot/database"
	"strconv"

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

	database.Initialize()

	statusMap, c, m := make(StatusMap), make(chan StatusUpdateCommand), make(chan StatusUpdateMonth)

	go manageStatus(statusMap, c, m)

	updates := bot.GetUpdatesChan(update_config)

	for update := range updates {
		if update.Message != nil {
			go handleUpdate(bot, update.Message, statusMap, c, m)
		}
	}
}

func manageStatus(
	status StatusMap,
	c <-chan StatusUpdateCommand,
	m <-chan StatusUpdateMonth,
) {
	for {
		select {
		case update := <-c:
			{
				if entry, ok := status[update.id]; ok {
					log.Printf("Updating: %s", update.command)
					// Then we modify the copy
					entry.command = update.command
				 
					// Then we reassign map entry
					status[update.id] = entry
				} else {
					status[update.id] = userStatusNew(update.command, 0, 0)
				}
	
				log.Printf("Command set, status: %+v", status)
			}
		case update := <-m:
			{
				if entry, ok := status[update.id]; ok {
					log.Printf("Updating: %d", update.month)
					// Then we modify the copy
					entry.month = update.month
				 
					// Then we reassign map entry
					status[update.id] = entry
				} else {
					status[update.id] = userStatusNew("", update.month, 0)
				}
	
				log.Printf("Month set, status: %+v", status)
			}
		}
	}

	// for update := range c {
	// 	status[update.id] = userStatusNew(update.command, 0, 0)

	// 	log.Printf("Status: %+v", status)
	// }

}

func handleUpdate(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	status StatusMap,
	c chan<- StatusUpdateCommand,
	m chan<- StatusUpdateMonth,
) {
	log.Printf("[%s] %s", message.From.UserName, message.Text)

	userId := message.From.ID

	c <- statusUpdateCommandNew(userId, message.Text)

	reply := tgbotapi.NewMessage(message.Chat.ID, message.Text)
	reply.ReplyToMessageID = message.MessageID

	if userStatus, exists := status[userId]; exists {
		if userStatus.command == "remember-month" {
			log.Printf("Last command was remember-month")

			month, _ := strconv.ParseUint(reply.Text, 10, 8)
			m <- statusUpdateMonthNew(userId, uint8(month))
			c <- statusUpdateCommandNew(userId, "remember-day")

			reply.Text = "Ok, che giorno?"
			reply.ReplyMarkup = dayKeyboard


			bot.Send(reply)
			return
		}
	} else {
		log.Printf("User ID %d still not present into status map", userId)
	}

	switch message.Text {
	case "month":
		reply.ReplyMarkup = monthKeyboard
	case "day":
		reply.ReplyMarkup = dayKeyboard
	case "last":
		{
			if userStatus, exists := status[userId]; exists {
				reply.Text = "Last command: " + userStatus.command
			} else {
				log.Printf("User ID %d still not present into status map", userId)
			}
		}
	case "close":
		reply.ReplyMarkup = emptyKeyboard
	case "ricorda":
		{
			log.Printf("Command is ricorda")
			c <- statusUpdateCommandNew(userId, "remember-month")
			
			reply.Text = "Ok, che mese?"
			reply.ReplyMarkup = monthKeyboard

			// database.UserInsert(userId, "Marcello")
			// database.BirthdayInsert("Mario", 27, 3, userId)
		}
	}
	

	bot.Send(reply)
}
