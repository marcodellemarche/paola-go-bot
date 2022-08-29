package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"paola-go-bot/database"
	"paola-go-bot/telegram/status"
	"paola-go-bot/telegram/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var WishlistsBirthdaysGet = Command{
	Name:        "regali",
	Description: "Lista dei regali dei miei amici",
	Handle:      handleWishlistsBirthdaysGet,
}

func handleWishlistsBirthdaysGet(message *tgbotapi.Message) status.CommandResponse {
	log.Printf("Get wishlists for my birthdays")

	birthdays, ok := database.BirthdayFind(0, 0, message.Chat.ID)
	if !ok {
		log.Printf("Error getting wishlists birthdays, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if len(birthdays) == 0 {
		log.Printf("Warning getting birthdays, no birthdays found yet")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Non ci sono regali ancora 🥲")
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	var rows [][]string
	var birthdayArgs []database.Birthday
	for _, birthday := range birthdays {
		if birthday.ContactId > 0 {
			rows = append(rows, []string{birthday.Name})
			birthdayArgs = append(birthdayArgs, birthday)
		}
	}
	wishlistsBirthdaysKeyboard := utils.Keyboard(rows...)

	if len(birthdayArgs) == 0 {
		log.Printf("Warning getting birthdays, no birthdays with contact id found yet")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Non ci sono regali ancora 🥲")
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	bytesBirthdayArgs, err := json.Marshal(birthdayArgs)
	if err != nil {
		log.Printf("Error getting wishlists birthdays, could not marshal birthday args into json")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, di chi?")
	status.SetNext(message.Chat.ID, askWhichWishlistToGet, string(bytesBirthdayArgs))
	return status.CommandResponse{Reply: &reply, Keyboard: &wishlistsBirthdaysKeyboard}
}

func askWhichWishlistToGet(message *tgbotapi.Message, args ...string) status.CommandResponse {
	log.Printf("Get wishlists - ask which wishlist to get for selected user")

	userName := message.Text

	var birthdayArgs []database.Birthday
	json.Unmarshal([]byte(args[0]), &birthdayArgs)

	var userId int64
	for _, birthdayArg := range birthdayArgs {
		if birthdayArg.Name == userName {
			userId = birthdayArg.ContactId
			break
		}
	}

	wishlists, ok := database.WishlistFindByUser(userId)
	if !ok {
		log.Printf("Error getting wishlists birthdays, could not fetch database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if len(wishlists) == 0 {
		log.Printf("Warning getting wishlists birthdays, no wishlists found yet")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Non ci sono regali ancora 🥲")
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	var rows [][]string
	for _, wishlist := range wishlists {
		rows = append(rows, []string{wishlist.Name})
	}
	wishlistsKeyboard := utils.Keyboard(rows...)

	bytesWishlistArgs, err := json.Marshal(wishlists)
	if err != nil {
		log.Printf("Error getting wishlists birthdays, could not marshal wishlist args into json")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, "Ok, quale regalo?")
	status.SetNext(message.Chat.ID, showWishlistDetails, userName, strconv.Itoa(int(userId)), string(bytesWishlistArgs))
	return status.CommandResponse{Reply: &reply, Keyboard: &wishlistsKeyboard}
}

func showWishlistDetails(message *tgbotapi.Message, args ...string) status.CommandResponse {
	log.Printf("Get wishlists - show wishlist details")

	wishlistName := message.Text

	userName := args[1]
	// userId, _ := strconv.ParseInt(args[2], 10, 64)

	var wishlistArgs []database.Wishlist
	json.Unmarshal([]byte(args[3]), &wishlistArgs)

	var wishlist database.Wishlist
	for _, wishlistArg := range wishlistArgs {
		if wishlistArg.Name == wishlistName {
			wishlist = wishlistArg
			break
		}
	}

	reply := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("%s vorrebbe '%s':\n", userName, wishlist.Name))

	if wishlist.Link != "" {
		reply.Text += fmt.Sprintf("\n%s\n", wishlist.Link)
	}

	if wishlist.BuyerId == message.Chat.ID {
		reply.Text += "\nLo hai già prenotato tu!"
		// TODO: cancella prenotazione
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	} else if wishlist.BuyerId > 0 {
		reply.Text += "\nGià prenotato da qualcun altro 🥲"
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	} else {
		reply.Text += "\nNon ancora prenotato! Vuoi regalarlo tu?"
		yesOrNoKeyboard := utils.Keyboard(
			[]string{"Sì"},
			[]string{"No"},
		)
		status.SetNext(message.Chat.ID, confirmWishlistBooking, wishlist.Name)
		return status.CommandResponse{Reply: &reply, Keyboard: &yesOrNoKeyboard}
	}

	// status.ResetNext(message.Chat.ID)
	// return status.CommandResponse{Reply: &reply, Keyboard: nil}
}

func confirmWishlistBooking(message *tgbotapi.Message, args ...string) status.CommandResponse {
	log.Printf("Get wishlists - confirm wishlist booking")

	userId, _ := strconv.ParseInt(args[2], 10, 64)
	wishlistName := args[4]

	var confirmed bool
	if message.Text == "Sì" {
		confirmed = true
	} else if message.Text == "No" {
		confirmed = false
	} else {
		log.Printf("Error booking wishlist, answer is not valid: %s", message.Text)
		reply := tgbotapi.NewMessage(message.Chat.ID, "Ma che risposta è, o sì o no, è facile")
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	if !confirmed {
		log.Printf("Error booking wishlist, booking not confirmed")
		reply := tgbotapi.NewMessage(message.Chat.ID, "Allora fa come cazzo te pare")
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	ok := database.WishlistSetBuyer(wishlistName, userId, message.Chat.ID)
	if !ok {
		log.Printf("Error booking wishlist, could not update database")
		reply := tgbotapi.NewMessage(message.Chat.ID, errorMessage)
		status.ResetNext(message.Chat.ID)
		return status.CommandResponse{Reply: &reply, Keyboard: nil}
	}

	log.Printf("Wishlist booking confirmed")
	reply := tgbotapi.NewMessage(message.Chat.ID, "Ottimo, lo farai tu allora!")
	status.ResetNext(message.Chat.ID)
	return status.CommandResponse{Reply: &reply, Keyboard: nil}
}
