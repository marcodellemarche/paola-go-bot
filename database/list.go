package database

import (
	"log"

	_ "embed"

	_ "github.com/lib/pq"
)

//go:embed list/create.sql
var sqlListCreateTable string

//go:embed list/drop.sql
var sqlListDropTable string

//go:embed list/insert.sql
var sqlListInsert string

//go:embed list/delete_by_subscriber.sql
var sqlListDeleteBySubscriber string

//go:embed list/find_by_user.sql
var sqlListFindByUser string

type List struct {
	UserId       int64
	SubscriberId int64
	UserName     string
}

func ListCreateTable() bool {
	_, err := db.Exec(sqlListCreateTable)
	if err != nil {
		log.Printf("Error creating list table: %s", err.Error())
		return false
	}

	return true
}

func ListDropTable() bool {
	_, err := db.Exec(sqlListDropTable)
	if err != nil {
		log.Printf("Error dropping list table: %s", err.Error())
		return false
	}

	return true
}

func ListInsert(userId int64, subscriberId int64, userName string) bool {
	_, err := db.Exec(sqlListInsert, userId, subscriberId, userName)
	if err != nil {
		log.Printf("Error inserting list into database, list insertion failed: %s", err.Error())
		return false
	}

	return true
}

func ListFindByUser(userId int64) (List, bool) {
	var list List

	rows, err := db.Query(sqlListFindByUser, userId)
	if err != nil {
		log.Printf("Error fetching list from database, query failed: %s", err.Error())
		return list, false
	}

	if !rows.Next() {
		log.Printf("No list for user %d found", userId)
		return list, true
	}

	defer rows.Close()

	rows.Scan(&list.UserId, &list.SubscriberId, &list.UserName)
	if err != nil {
		log.Printf("Error fetching list from database, scan failed: %s", err.Error())
		return list, false
	}

	return list, true
}

func ListDeleteBySubscriber(subscriberId int64, userId int64) bool {
	_, err := db.Exec(sqlListDeleteBySubscriber, subscriberId, userId)
	if err != nil {
		log.Printf("Error deleting list from database, list deletion failed: %s", err.Error())
		return false
	}

	return true
}
