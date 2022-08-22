package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Keyboard(rows ...[]string) tgbotapi.ReplyKeyboardMarkup {
	var keyboardRows [][]tgbotapi.KeyboardButton = make([][]tgbotapi.KeyboardButton, len(rows))

	for r, row := range rows {
		var keyboardRow []tgbotapi.KeyboardButton = make([]tgbotapi.KeyboardButton, len(row))

		for c, column := range row {
			keyboardRow[c] = tgbotapi.NewKeyboardButton(column)
		}

		keyboardRows[r] = tgbotapi.NewKeyboardButtonRow(keyboardRow...)
	}

	return tgbotapi.NewReplyKeyboard(keyboardRows...)
}

var EmptyKeyboard = tgbotapi.NewRemoveKeyboard(true)

var MonthKeyboard = Keyboard(
	[]string{"1", "2", "3"},
	[]string{"4", "5", "6"},
	[]string{"7", "8", "9"},
	[]string{"10", "11", "12"},
)

var DayKeyboard = Keyboard(
	[]string{"1", "2", "3", "4", "5", "6", "7"},
	[]string{"8", "9", "10", "11", "12", "13"},
	[]string{"14", "15", "16", "17", "18", "19"},
	[]string{"20", "21", "22", "23", "24", "25"},
	[]string{"26", "27", "28", "29", "30", "31"},
)
