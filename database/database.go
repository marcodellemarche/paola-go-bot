package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "ec2-52-212-228-71.eu-west-1.compute.amazonaws.com"
	port     = 5432
	user     = "xuwhqoruiguhbd"
	password = "9368fef0eee2b6b71a5cab6bce7ba4879b3124130686ccdef5973917695e8f8d"
	dbname   = "dd6snkdliqfa7q"
)

var db sql.DB

func Initialize() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)

	pDb, err := sql.Open("postgres", psqlconn)
	db = *pDb
	CheckError(err)

	// defer db.Close()

	err = db.Ping()
	CheckError(err)

	log.Println("DB connected!")

	UserCreateTable()
	BirthdayCreateTable()
}

func CheckError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
