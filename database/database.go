package database

import (
	"log"

	"database/sql"

	_ "github.com/lib/pq"
)

var db sql.DB

func Initialize(psqlconn string, debug bool) {
	// psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)
	pDb, err := sql.Open("postgres", psqlconn)
	db = *pDb
	if err != nil {
		log.Fatalf("Error opening database: %s", err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %s", err.Error())
	}

	log.Println("DB connected!")

	UserCreateTable()
	BirthdayCreateTable()
	ListCreateTable()
}
