package database

import (
	"log"
	"strings"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

var (
	userCollection = f.Collection("users")
)

type User struct {
	id    string
	Name  string `fauna:"name"`
	Super bool   `fauna:"super"`
}

func userNew(id string, name string, super bool) User {
	return User{
		id,
		name,
		super,
	}
}

func createUserCollection(client *f.FaunaClient) {
	_, err := client.Query(f.CreateCollection(f.Obj{"name": "users"}))

	if err != nil && !strings.Contains(err.Error(), "instance already exists") {
		log.Fatalln(err)
	}
}

func createUser(client *f.FaunaClient, user User) string {
	var userId f.RefV

	newUser, err := client.Query(
		f.Create(
			f.Ref(userCollection, user.id),
			f.Obj{"data": user},
		),
	)

	if err != nil {
		log.Fatalln(err)
	}

	err = newUser.At(ref).Get(&userId)

	if err != nil {
		log.Fatalln(err)
	}

	return userId.ID
}

func updateUser(client *f.FaunaClient, id string, user User) {
	_, err := client.Query(
		f.Update(
			f.Ref(userCollection, id),
			f.Obj{"data": user},
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func getUser(client *f.FaunaClient, id string) User {
	var user User

	value, err := client.Query(f.Get(f.Ref(userCollection, id)))

	if err != nil {
		log.Fatalln(err)
	}

	err = value.At(data).Get(&user)

	if err != nil {
		log.Fatalln(err)
	}

	user.id = id

	return user
}

func deleteUser(client *f.FaunaClient, id string) {
	_, err := client.Query(f.Delete(f.Ref(userCollection, id)))

	if err != nil {
		log.Fatalln(err)
	}
}
