package scripts

import (
	_ "embed"
	"log"
	"os"
	"time"

	"encoding/json"

	"paola-go-bot/database"
	"paola-go-bot/telegram"
)

//go:embed chats.json
var jsonMigrationChats string

//go:embed super.json
var jsonMigrationSuper string

type User struct {
	// _id      string `json:"_id"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
}

type Chat struct {
	// _id    string `json:"_id"`
	ChatId int64  `json:"chatId"`
	Name   string `json:"name"`
	Trust  bool   `json:"trust"`
	Users  []User `json:"users"`
	// __v    int64  `json:"__v"`
}

var SuperPaolaId int64 = 302635332

func Migration(token string, write bool) {
	if token == "" {
		token = os.Getenv("TELEGRAM_TOKEN")
	}

	if write {
		database_uri := os.Getenv("DATABASE_URI")

		database.Initialize(database_uri, false)

		database.ListDropTable()
		database.BirthdayDropTable()
		database.UserDropTable()

		database.UserCreateTable()
		database.BirthdayCreateTable()
		database.ListCreateTable()
	}

	telegram.Initialize(token, false)

	database.UserInsert(SuperPaolaId, "SuperPaola")

	var chats []Chat
	json.Unmarshal([]byte(jsonMigrationChats), &chats)

	for i, chat := range chats {
		log.Printf("Chat %d: %d - %s - super %v", i, chat.ChatId, chat.Name, chat.Trust)

		name := telegram.GetNameFromUserId(chat.ChatId)

		log.Printf("- name: %s", name)

		if write {
			database.UserInsert(chat.ChatId, chat.Name)
		}

		if chat.Trust {
			if write {
				database.ListInsert(SuperPaolaId, chat.ChatId)
			}
		}

		for j, user := range chat.Users {
			log.Printf("    User %d: %s - %s", j, user.Name, user.Birthday)

			birthday, err := time.Parse("2006-01-02T00:00:00.000Z", user.Birthday)
			if err != nil {
				log.Printf("Error during migration, parsing date %s failed: %s", user.Birthday, err.Error())
				return
			}

			if write {
				database.BirthdayInsert(user.Name, 0, uint8(birthday.Day()), uint8(birthday.Month()), chat.ChatId)
			}
		}
	}

	var super Chat
	json.Unmarshal([]byte(jsonMigrationSuper), &super)

	log.Printf("SuperPaola: %d", SuperPaolaId)

	for j, user := range super.Users {
		log.Printf("    User %d: %s - %s", j, user.Name, user.Birthday)

		birthday, err := time.Parse("2006-01-02T00:00:00.000Z", user.Birthday)
		if err != nil {
			log.Printf("Error during migration, parsing date %s failed: %s", user.Birthday, err.Error())
			return
		}

		if write {
			database.BirthdayInsert(user.Name, 0, uint8(birthday.Day()), uint8(birthday.Month()), SuperPaolaId)
		}
	}
}
