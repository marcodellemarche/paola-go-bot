package commands

import (
	"log"
	"net/url"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"
	"paola-go-bot/telegram/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var WishlistSet = Command{
	Name:        "desidera",
	Description: "Desidera un regalo",
	Handle:      handleWishlistSet,
}

func handleWishlistSet(message *tgbotapi.Message) status.CommandResponse {
	log.Printf("Set wishlist - asking for name")

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, come si chiama il regalo?")
	status.SetNext(message.Chat.ID, askForWishlistLink)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}
}

func askForWishlistLink(message *tgbotapi.Message, args ...string) status.CommandResponse {
	log.Printf("Set wishlist - received name, asking for link")

	name := message.Text
	
	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, hai un link?")
	status.SetNext(message.Chat.ID, confirmWishlist, name)
	noLinkKeyboard := utils.Keyboard([]string{"No"})
	return status.CommandResponse{Reply: &reply, Keyboard: &noLinkKeyboard}
}

func confirmWishlist(message *tgbotapi.Message, args ...string) status.CommandResponse {
	log.Printf("Set wishlist - received link, confirming wishlist")

	name := args[0]
	link := message.Text

	if name == "" {
		log.Printf("Error confirming wishlist, name is not valid: <empty-string>")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il nome non è valido!")
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if link == "No" {
		link = ""
	} else {
		_, err := url.ParseRequestURI(link)
		if err != nil {
			log.Printf("Error confirming wishlist, link is not valid: <invalid-url>")
			reply := tgbotapi.NewMessage(message.Chat.ID, "Oh, ma il link non è valido!")
			status.ResetNext(message.Chat.ID)
			return status.CommandResponse{Reply: &reply, Keyboard: nil}
		}
	}

	ok := database.WishlistInsert(message.Chat.ID, name, link, 0)
	if !ok {
		log.Printf("Error confirming wishlist, could not update database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, regalo ricevuto ✌️")
	status.ResetNext(message.Chat.ID)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}
}
