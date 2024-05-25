package telegram

import (
	"log"

	"paola-go-bot/telegram/commands"
	"paola-go-bot/telegram/status"
	"paola-go-bot/telegram/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Telegram) aiHandleText(message *tgbotapi.Message) status.CommandResponse {
	log.Printf("Generic text command - asking to PaolaGPT")

	if !t.chatgptClient.RateLimiter.Allow(message.Chat.ID) {
		log.Println("Rate limit exceeded")

		reply := tgbotapi.NewMessage(message.Chat.ID, utils.RandomInsult())
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	var threadID string
	if userStatus, exists := status.Get(message.Chat.ID); exists {
		threadID = userStatus.ThreadID
	}

	answer, err := t.useAssistant(message.Text, message.Chat.ID, &threadID)
	if err != nil {
		log.Printf("Error using the AI assistant: %s", err)

		reply := tgbotapi.NewMessage(message.Chat.ID, commands.ErrorMessage)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if threadID != "" {
		status.SetThread(message.Chat.ID, threadID)
	} else {
		status.ResetThread(message.Chat.ID)
	}

	status.ResetNext(message.Chat.ID)
	reply := tgbotapi.NewMessage(message.Chat.ID, answer)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}

}
