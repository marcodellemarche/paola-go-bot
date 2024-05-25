package status

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NextCommand func(message *tgbotapi.Message, args ...string) CommandResponse

type CommandResponse struct {
	Reply    *tgbotapi.MessageConfig
	Keyboard *tgbotapi.ReplyKeyboardMarkup
}

type TelegramStatus struct {
	Next     NextCommand
	Args     []string
	ThreadID string
}

type UserId = int64
