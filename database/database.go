package database

import (
	"fmt"
	"log"

	"database/sql"

	_ "github.com/lib/pq"
)

var db sql.DB

var host string = "database"
var port int16 = 5432
var sslmode string = "disable" // Use "require" for cloud DBs and "disable" for self hosted ones

func Initialize(user string, dbname string, password string, debug bool) {
	if user == "" {
		log.Fatal("Missing database user")
	}

	if dbname == "" {
		log.Fatal("Missing database name")
	}

	if password == "" {
		log.Fatal("Missing database password")
	}

	databaseUri := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	pDb, err := sql.Open("postgres", databaseUri)
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
