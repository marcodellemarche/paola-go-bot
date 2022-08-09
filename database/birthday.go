package database

import (
	"log"
	"strings"

	f "github.com/fauna/faunadb-go/v4/faunadb"
)

var (
	birthdayCollection = f.Collection("birthdays")
)

type Birthday struct {
	Name string `fauna:"name"`
	Date int64  `fauna:"date"`
}

func birthdayNew(name string, date int64) Birthday {
	return Birthday{
		name,
		date,
	}
}

func createBirthdayCollection(client *f.FaunaClient) {
	_, err := client.Query(f.CreateCollection(f.Obj{"name": "birthdays"}))

	if err != nil && !strings.Contains(err.Error(), "instance already exists") {
		log.Fatalln(err)
	}
}

func createBirthday(client *f.FaunaClient, birthday Birthday) string {
	var birthdayId f.RefV

	newBirthday, err := client.Query(
		f.Create(
			f.Collection(birthdayCollection),
			f.Obj{"data": birthday},
		),
	)

	if err != nil {
		log.Fatalln(err)
	}

	err = newBirthday.At(ref).Get(&birthdayId)

	if err != nil {
		log.Fatalln(err)
	}

	return birthdayId.ID
}

func updateBirthday(client *f.FaunaClient, id string, birthday Birthday) {
	_, err := client.Query(
		f.Update(
			f.Ref(birthdayCollection, id),
			f.Obj{"data": birthday},
		),
	)

	if err != nil {
		log.Fatalln(err)
	}
}

func getBirthday(client *f.FaunaClient, id string) Birthday {
	var birthday Birthday

	value, err := client.Query(f.Get(f.Ref(birthdayCollection, id)))

	if err != nil {
		log.Fatalln(err)
	}

	err = value.At(data).Get(&birthday)

	if err != nil {
		log.Fatalln(err)
	}

	return birthday
}

func deleteBirthday(client *f.FaunaClient, id string) {
	_, err := client.Query(f.Delete(f.Ref(birthdayCollection, id)))

	if err != nil {
		log.Fatalln(err)
	}
}
