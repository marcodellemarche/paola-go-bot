package database

import (
	"log"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

var (
	data = f.ObjKey("data")
	ref  = f.ObjKey("ref")
)

var client *f.FaunaClient

func Initialize(secret string) {
	client = f.NewFaunaClient(secret)

	createUserCollection(client)
	createBirthdayCollection(client)
}

func Start() {
	if client == nil {
		log.Fatalln("Client not initialized")
	}

	userId := "1234"

	// birthdays := []Birthday{
	// 	birthdayNew("pippo", time.Date(2000, time.January, 10, 0, 0, 0, 0, time.UTC).Unix()),
	// 	birthdayNew("franco", time.Date(2000, time.November, 16, 0, 0, 0, 0, time.UTC).Unix()),
	// }

	user := userNew(userId, "mario", false)

	log.Printf("User: %v", user)

	_ = createUser(client, user)

	log.Printf("User ID: %s", userId)

	user = getUser(client, userId)

	log.Printf("Created user: %v", user)

	user.Super = true
	updateUser(client, userId, user)

	user = getUser(client, userId)

	log.Printf("Updated user: %v", user)

	// deleteUser(client, userId)
}
