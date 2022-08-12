package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserId = int64

// type Next func(
// 	bot *tgbotapi.BotAPI,
// 	message *tgbotapi.Message,
// 	status StatusMap,
// 	c chan<- StatusUpdate,
// 	args ...string,
// )

type UserStatus struct {
	next func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, c chan<- StatusUpdate, args ...string)
	args []string
}

func UserStatusNew(next func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, c chan<- StatusUpdate, args ...string), args ...string) UserStatus {
	return UserStatus{
		next,
		args,
	}
}

type StatusMap = map[UserId]UserStatus

type StatusUpdate struct {
	id   UserId
	next func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, c chan<- StatusUpdate, args ...string)
	args []string
}

func StatusUpdateNew(id UserId, next func(bot *tgbotapi.BotAPI, message *tgbotapi.Message, c chan<- StatusUpdate, args ...string), args ...string) StatusUpdate {
	return StatusUpdate{
		id,
		next,
		args,
	}
}
