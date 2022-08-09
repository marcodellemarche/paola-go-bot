package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func keyboard(rows ...[]string) tgbotapi.ReplyKeyboardMarkup {
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
