package database

import (
	_ "embed"
	"fmt"
	"log"
	"sort"
	"time"

	"database/sql"

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

//go:embed birthday/find_by_date.sql
var sqlBirthdayFindByDate string

//go:embed birthday/find_by_date_and_list.sql
var sqlBirthdayFindByDateAndList string

//go:embed birthday/delete_by_name.sql
var sqlBirthdayDeleteByName string

type Birthday struct {
	Name      string
	Day       uint8
	Month     uint8
	date      time.Time
	contactId int64
	UserId    int64
}

func BirthdayCreateTable() bool {
	_, err := db.Exec(sqlBirthdayCreateTable)
	if err != nil {
		log.Printf("Error creating birthday table: %s", err.Error())
		return false
	}

	return true
}

func BirthdayDropTable() bool {
	_, err := db.Exec(sqlBirthdayDropTable)
	if err != nil {
		log.Printf("Error dropping birthday table: %s", err.Error())
		return false
	}

	return true
}

func BirthdayInsert(name string, contactId int64, day uint8, month uint8, userId int64) bool {
	if day > 31 {
		log.Printf("Error inserting birthday into database, invalid day: %d > 31", day)
		return false
	}

	if month > 12 {
		log.Printf("Error inserting birthday into database, invalid month: %d > 12", month)
		return false
	}

	formattedDate := fmt.Sprintf("2000-%02d-%02d", month, day)

	date, err := time.Parse("2006-01-02", formattedDate)
	if err != nil {
		log.Printf("Error inserting birthday into database, date format is not valid: %s - %s", formattedDate, err.Error())
		return false
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error inserting birthday into database, transaction begin failed: %s", err.Error())
		return false
	}

	nullableContactId := sql.NullInt64{Int64: contactId, Valid: contactId > 0}

	if nullableContactId.Valid {
		_, err = tx.Exec(sqlUserInsert, contactId, name)
		if err != nil {
			log.Printf("Error inserting birthday into database, contact insertion failed: %s", err.Error())
			_ = tx.Rollback()
			return false
		}
	}

	_, err = tx.Exec(sqlBirthdayInsert, name, nullableContactId, date, userId)
	if err != nil {
		log.Printf("Error inserting birthday into database, birthday insertion failed: %s", err.Error())
		_ = tx.Rollback()
		return false
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Error inserting birthday into database, transaction commit failed: %s", err.Error())
		_ = tx.Rollback()
		return false
	}

	return true
}

func BirthdayFindByUser(userId int64) ([]Birthday, bool) {
	rows, err := db.Query(sqlBirthdayFindByUser, userId)
	if err != nil {
		log.Printf("Error finding birthdays from database, query failed: %s", err.Error())
		return nil, false
	}

	return birthdayFind(rows, true)
}

func BirthdayFindByDate(day uint8, month uint8) ([]Birthday, bool) {
	rows, err := db.Query(sqlBirthdayFindByDate, day, month)
	if err != nil {
		log.Printf("Error finding birthdays from database, query failed: %s", err.Error())
		return nil, false
	}

	return birthdayFind(rows, false)
}

func BirthdayFindByDateAndList(day uint8, month uint8, listId int64, subscriberId int64) ([]Birthday, bool) {
	nullableSubscriberId := sql.NullInt64{Int64: subscriberId, Valid: subscriberId > 0}

	rows, err := db.Query(sqlBirthdayFindByDateAndList, day, month, listId, nullableSubscriberId)
	if err != nil {
		log.Printf("Error finding birthdays from database, query failed: %s", err.Error())
		return nil, false
	}

	return birthdayFind(rows, false)
}

func birthdayFind(rows *sql.Rows, toSort bool) ([]Birthday, bool) {
	birthdays := make([]Birthday, 0)
	for rows.Next() {
		var birthday Birthday
		var formattedDate string
		var contactId sql.NullInt64

		err := rows.Scan(&birthday.Name, &contactId, &formattedDate, &birthday.UserId)
		if err != nil {
			log.Printf("Error finding birthdays from database, scan failed: %s", err.Error())
			return nil, false
		}

		birthday.date, err = time.Parse("2006-01-02T00:00:00Z", formattedDate)
		if err != nil {
			log.Printf("Error finding birthdays from database, parsing date %s failed: %s", formattedDate, err.Error())
			return nil, false
		}

		birthday.Day = uint8(birthday.date.Day())
		birthday.Month = uint8(birthday.date.Month())

		if contactId.Valid {
			birthday.contactId = contactId.Int64
		}

		birthdays = append(birthdays, birthday)
	}

	if toSort {
		sort.Slice(birthdays, func(i, j int) bool {
			return birthdays[i].date.Before(birthdays[j].date)
		})
	}

	return birthdays, true
}

func BirthdayDeleteByName(name string, userId int64) bool {
	_, err := db.Exec(sqlBirthdayDeleteByName, name, userId)
	if err != nil {
		log.Printf("Error deleting birthday by name from database, deletion failed: %s", err.Error())
		return false
	}

	return true
}
