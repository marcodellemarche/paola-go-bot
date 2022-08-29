package database

import (
	_ "embed"
	"fmt"
	"log"
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

//go:embed birthday/find.sql
var sqlBirthdayFind string

//go:embed birthday/find_by_list.sql
var sqlBirthdayFindByList string

//go:embed birthday/delete_by_name.sql
var sqlBirthdayDeleteByName string

var BaseYear = 2000

type Birthday struct {
	Name      string
	Day       uint8
	Month     uint8
	date      time.Time
	ContactId int64
	UserId    int64
	ListId    int64
	ListName  string
}

func (b *Birthday) Passed() bool {
	parsedNow := time.Now().AddDate(BaseYear - time.Now().Year(), 0, 0)

	return b.date.Before(parsedNow)
}

func (b *Birthday) Before(a *Birthday) bool {
	return b.date.Before(a.date)
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

	formattedDate := fmt.Sprintf("%d-%02d-%02d", BaseYear, month, day)

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

func BirthdayFind(day uint8, month uint8, userId int64) ([]Birthday, bool) {
	nullableDay := sql.NullInt64{Int64: int64(day), Valid: day > 0}
	nullableMonth := sql.NullInt64{Int64: int64(month), Valid: month > 0}
	nullableUserId := sql.NullInt64{Int64: userId, Valid: userId > 0}

	rows, err := db.Query(sqlBirthdayFind, nullableDay, nullableMonth, nullableUserId)
	if err != nil {
		log.Printf("Error finding birthdays from database, query failed: %s", err.Error())
		return nil, false
	}

	return birthdayFind(rows)
}

func BirthdayFindByList(day uint8, month uint8, listId int64, subscriberId int64) ([]Birthday, bool) {
	nullableDay := sql.NullInt64{Int64: int64(day), Valid: day > 0}
	nullableMonth := sql.NullInt64{Int64: int64(month), Valid: month > 0}
	nullableListId := sql.NullInt64{Int64: listId, Valid: listId > 0}
	nullableSubscriberId := sql.NullInt64{Int64: subscriberId, Valid: subscriberId > 0}

	rows, err := db.Query(sqlBirthdayFindByList, nullableDay, nullableMonth, nullableListId, nullableSubscriberId)
	if err != nil {
		log.Printf("Error finding birthdays from database, query failed: %s", err.Error())
		return nil, false
	}

	return birthdayFind(rows)
}

func birthdayFind(rows *sql.Rows) ([]Birthday, bool) {
	birthdays := make([]Birthday, 0)
	for rows.Next() {
		var birthday Birthday
		var formattedDate string
		var contactId sql.NullInt64
		var listId sql.NullInt64
		var listName sql.NullString

		err := rows.Scan(&birthday.Name, &contactId, &formattedDate, &birthday.UserId, &listId, &listName)
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
			birthday.ContactId = contactId.Int64
		}

		if listId.Valid {
			birthday.ListId = listId.Int64
		}

		if listName.Valid {
			birthday.ListName = listName.String
		}

		birthdays = append(birthdays, birthday)
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
