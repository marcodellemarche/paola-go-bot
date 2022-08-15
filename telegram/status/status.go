package status

import (
	"log"
)

var status map[UserId]TelegramStatus

var c chan Update

func init() {
	status = make(map[UserId]TelegramStatus)

	c = make(chan Update)
}

func Manage() {
	for update := range c {
		var args []string = update.Args

		if oldStatus, exists := status[update.Id]; exists && oldStatus.Next != nil {
			args = append(oldStatus.Args, args...)
		}

		status[update.Id] = TelegramStatus{
			update.Next,
			args,
		}

		log.Printf("Status updated for %d: %+v", update.Id, status[update.Id])
	}
}

func Get(userId int64) (TelegramStatus, bool) {
	status, exists := status[userId]

	return status, exists
}
