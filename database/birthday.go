package database

import (
	"fmt"
	"log"
	"time"

	_ "embed"

	_ "github.com/lib/pq"
)

//go:embed birthday/create.sql
var sqlBirthdayCreateTable string

//go:embed birthday/drop.sql
var sqlBirthdayDropTable string

//go:embed birthday/insert.sql
var sqlBirthdayInsert string

//go:embed birthday/find_by_user.sql
var sqlBirthdayFindByUser string

//go:embed birthday/delete_by_name.sql
var sqlBirthdayDeleteByName string

type Birthday struct {
	id      string
	name    string
	date    string
	user_id int64
}

func BirthdayNew(id string, name string, date string, user_id int64) Birthday {
	return Birthday{
		id,
		name,
		date,
		user_id,
	}
}

func BirthdayCreateTable() {
	_, err := db.Exec(sqlBirthdayCreateTable)
	CheckError(err)
}

func BirthdayDropTable() {
	_, err := db.Exec(sqlBirthdayDropTable)
	CheckError(err)
}

func BirthdayInsert(name string, day uint8, month uint8, userId int64) bool {
	if day > 31 {
		log.Printf("Error inserting birthday into database, invalid day: %d > 31", day);
		return false
	}

	if month > 12 {
		log.Printf("Error inserting birthday into database, invalid month: %d > 12", month);
		return false
	}

	formattedDate := fmt.Sprintf("2000-%02d-%02d", month, day)
	
	date, err := time.Parse("2006-01-02", formattedDate)
	if err != nil {
		log.Printf("Error inserting birthday into database, date format is not valid: %s", formattedDate)
		return false
	}

	_, err = db.Exec(sqlBirthdayInsert, name, date, userId)
	if err != nil {
		log.Printf("Error inserting birthday into database, generic error. Maybe the name %s is a duplicate?", name)
		return false
	}

	return true
}

func BirthdayFindByUser(userId int64) []Birthday {
	rows, err := db.Query(sqlBirthdayFindByUser, userId)
	CheckError(err)

	// defer rows.Close()

	birthdays := make([]Birthday, 0)
	for rows.Next() {
		var birthday Birthday
		rows.Scan(&birthday.id, &birthday.name, &birthday.date, &birthday.user_id)
		CheckError(err)
		birthdays = append(birthdays, birthday)
	}

	return birthdays
}

func BirthdayDeleteByName(name string) {
	_, err := db.Exec(sqlBirthdayDeleteByName, name)
	CheckError(err)
}
