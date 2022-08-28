package commands

import (
	"fmt"
	"log"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"
	"paola-go-bot/telegram/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var errorMessage = "So 'ncazzo io, ma qualcosa è andato storto 🥲"

type Command struct {
	Name        string
	Description string
	Handle      func(message *tgbotapi.Message) status.CommandResponse
}

func CheckIfNewUser(message *tgbotapi.Message, userId int64) status.CommandResponse {
	user, ok := database.UserFindById(message.Chat.ID)
	if !ok {
		log.Printf("Error finding user, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if user.Id == 0 {
		reply := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Te sei nuovo, inizia un po' con /%s", Start.Name))
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if _, exists := status.Get(userId); !exists {
		status.ResetNext(message.Chat.ID)
	}

	return status.CommandResponse{Reply: nil, Keyboard: nil}
}

func CheckNextActionOrDefault(message *tgbotapi.Message, userId int64) status.CommandResponse {
	if userStatus, exists := status.Get(userId); exists {
		if userStatus.Next != nil {
			return userStatus.Next(message, userStatus.Args...)
		} else {
			log.Printf("User ID %d has next callback set to nil", userId)
		}
	} else {
		log.Printf("User ID %d still not present into status map", userId)
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, utils.RandomInsult())
	status.ResetNext(message.Chat.ID)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}
}

var CommandsEnabled = []tgbotapi.BotCommand{
	{
		Command:     Start.Name,
		Description: Start.Description,
	},
	{
		Command:     BirthdaySet.Name,
		Description: BirthdaySet.Description,
	},
	{
		Command:     BirthdayDelete.Name,
		Description: BirthdayDelete.Description,
	},
	{
		Command:     BirthdaysGet.Name,
		Description: BirthdaysGet.Description,
	},
	// {
	// 	Command: SubscribeList.Name,
	// 	Description: SubscribeList.Description,
	// },
	{
		Command:     Stop.Name,
		Description: Stop.Description,
	},
	{
		Command:     WishlistSet.Name,
		Description: WishlistSet.Description,
	},
}
