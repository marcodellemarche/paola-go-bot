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

func AskForName(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	status StatusMap,
	c chan<- StatusUpdate,
	args ...string,
) {
	log.Printf("Set birthday - asking for name")

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, che nome?")
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = emptyKeyboard

	c <- StatusUpdateNew(message.From.ID, askForMonth)

	bot.Send(reply)
}

func askForMonth(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	status StatusMap,
	c chan<- StatusUpdate,
	args ...string,
) {
	log.Printf("Set birthday - received name, asking for month")

	name := message.Text

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, che mese?")
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = monthKeyboard

	c <- StatusUpdateNew(message.From.ID, askForDay, name)

	bot.Send(reply)
}

func askForDay(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	status StatusMap,
	c chan<- StatusUpdate,
	args ...string,
) {
	log.Printf("Set birthday - received month, asking for day")

	name := args[0]
	month := message.Text

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, che giorno?")
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = dayKeyboard

	c <- StatusUpdateNew(message.From.ID, confirmBirthday, name, month)

	bot.Send(reply)
}

func confirmBirthday(
	bot *tgbotapi.BotAPI,
	message *tgbotapi.Message,
	status StatusMap,
	c chan<- StatusUpdate,
	args ...string,
) {
	log.Printf("Set birthday - received day, confirming birthday")

	name := args[0]
	month := args[1]
	day := message.Text

	if name == "" {
		log.Printf("Error confirming birthday, name is not valid")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il nome non è valido!")
		bot.Send(reply)
		return
	}

	parsedMonth, err := strconv.ParseUint(month, 10, 8)
	if err != nil {
		log.Printf("Error confirming birthday, month is not valid")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il mese non è valido!")
		bot.Send(reply)
		return
	}

	parsedDay, err := strconv.ParseUint(day, 10, 8)
	if err != nil {
		log.Printf("Error confirming birthday, day is not valid")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il giorno non è valido!")
		bot.Send(reply)
		return
	}

	inserted := database.BirthdayInsert(name, uint8(parsedDay), uint8(parsedMonth), message.From.ID)
	if !inserted {
		log.Printf("Error confirming birthday, could not update database")
		reply := tgbotapi.NewMessage(message.Chat.ID, "So 'ncazzo io, ma qualcosa è andato storto!")
		bot.Send(reply)
		return
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Perfetto, me lo ricorderò!")
	reply.ReplyToMessageID = message.MessageID
	reply.ReplyMarkup = emptyKeyboard

	c <- StatusUpdateNew(message.From.ID, nil)

	bot.Send(reply)
}
