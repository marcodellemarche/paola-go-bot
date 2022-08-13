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
	CheckError(err)

	err = db.Ping()
	CheckError(err)

	log.Println("DB connected!")

	if debug {
		BirthdayDropTable()
		UserDropTable()
	}

	UserCreateTable()
	BirthdayCreateTable()
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
