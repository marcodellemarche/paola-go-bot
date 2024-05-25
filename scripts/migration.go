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

type User struct {
	// _id      string `json:"_id"`
	Name     string `json:"name"`
	Birthday string `json:"birthday"`
}

type Chat struct {
	// _id    string `json:"_id"`
	ChatId   int64  `json:"chatId"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Trust    bool   `json:"trust"`
	Users    []User `json:"users"`
	// __v    int64  `json:"__v"`
}

var SuperPaolaId int64 = 302635332

var SuperPaolaName string = "Paola"

func Migration(telegram_token string, write bool) {
	if telegram_token == "" {
		telegram_token = os.Getenv("TELEGRAM_TOKEN")
	}

	jsonMigrationChats, err := os.ReadFile("chats.json")
	if err != nil {
		log.Fatalf("Error loading 'chats.json' file: %s", err.Error())
	}

	jsonMigrationSuper, err := os.ReadFile("super.json")
	if err != nil {
		log.Fatalf("Error loading 'super.json' file: %s", err.Error())
	}

	if write {
		database_user := os.Getenv("POSTGRES_USER")
		database_db := os.Getenv("POSTGRES_DB")
		database_password := os.Getenv("POSTGRES_PASSWORD")

		database.Initialize(database_user, database_db, database_password, false)

		database.ListDropTable()
		database.BirthdayDropTable()
		database.UserDropTable()

		database.UserCreateTable()
		database.BirthdayCreateTable()
		database.ListCreateTable()
	}

	telegramBot := telegram.New(telegram_token, nil, false)

	if write {
		database.UserInsert(SuperPaolaId, SuperPaolaName)
	}

	var chats []Chat
	json.Unmarshal(jsonMigrationChats, &chats)

	for i, chat := range chats {
		log.Printf("Chat %d: %d - %s - super %v", i, chat.ChatId, chat.Name, chat.Trust)

		name := telegramBot.GetNameFromUserId(chat.ChatId)

		log.Printf("- name: %s", name)

		if write {
			database.UserInsert(chat.ChatId, chat.Name)
		}

		if chat.Trust {
			if write {
				database.ListInsert(SuperPaolaId, chat.ChatId, SuperPaolaName)
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
