package status

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type NextFunc func(message *tgbotapi.Message, args ...string)

type TelegramStatus struct {
	Next NextFunc
	Args []string
}

type UserId = int64
