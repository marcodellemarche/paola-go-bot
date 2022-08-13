package telegram

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserId = int64

// type Next func(
// 	message *tgbotapi.Message,
// 	status StatusMap,
// 	c chan<- StatusUpdate,
// 	args ...string,
// )

type UserStatus struct {
	next func(message *tgbotapi.Message, c chan<- StatusUpdate, args ...string)
	args []string
}

func UserStatusNew(next func(message *tgbotapi.Message, c chan<- StatusUpdate, args ...string), args ...string) UserStatus {
	return UserStatus{
		next,
		args,
	}
}

type StatusMap = map[UserId]UserStatus

type StatusUpdate struct {
	id   UserId
	next func(message *tgbotapi.Message, c chan<- StatusUpdate, args ...string)
	args []string
}

func StatusUpdateNew(id UserId, next func(message *tgbotapi.Message, c chan<- StatusUpdate, args ...string), args ...string) StatusUpdate {
	return StatusUpdate{
		id,
		next,
		args,
	}
}

func manageStatus(
	status StatusMap,
	c <-chan StatusUpdate,
) {
	for update := range c {
		status[update.id] = UserStatusNew(update.next, update.args...)

		log.Printf("Status updated for %d: %+v", update.id, status[update.id])
	}
}
