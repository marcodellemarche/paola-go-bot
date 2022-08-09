package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "embed"

	_ "github.com/lib/pq"
)

const (
	host     = "ec2-52-212-228-71.eu-west-1.compute.amazonaws.com"
	port     = 5432
	user     = "xuwhqoruiguhbd"
	password = "9368fef0eee2b6b71a5cab6bce7ba4879b3124130686ccdef5973917695e8f8d"
	dbname   = "dd6snkdliqfa7q"
)

//go:embed user/create.sql
var sqlUserCreate string

//go:embed user/drop.sql
var sqlUserDrop string

//go:embed user/insert.sql
var sqlUserInsert string

//go:embed user/delete_by_id.sql
var sqlUserDeleteById string

//go:embed birthday/create.sql
var sqlBirthdayCreate string

//go:embed birthday/drop.sql
var sqlBirthdayDrop string

//go:embed birthday/insert.sql
var sqlBirthdayInsert string

//go:embed birthday/find_by_user.sql
var sqlBirthdayFindByUser string

//go:embed birthday/delete_by_name.sql
var sqlBirthdayDeleteByName string

func Sql() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	// defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	log.Println("Connected!")

	_, err = db.Exec(sqlUserCreate)
	CheckError(err)

	_, err = db.Exec(sqlBirthdayCreate)
	CheckError(err)

	_, err = db.Exec(sqlBirthdayDeleteByName, "Luca")
	CheckError(err)

	_, err = db.Exec(sqlUserDeleteById, 1234)
	CheckError(err)

	_, err = db.Exec(sqlUserInsert, 1234, "Mario")
	CheckError(err)

	_, err = db.Exec(sqlBirthdayInsert, "Luca", "2022-03-08", 1234)
	CheckError(err)

	rows, err := db.Query(sqlBirthdayFindByUser, 1234)
	CheckError(err)

	// defer rows.Close()
	for rows.Next() {
		var id string
		var name string
		var date string
		var user_id int64

		err = rows.Scan(&id, &name, &date, &user_id)
		CheckError(err)

		log.Println(id, name, date, user_id)
	}
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
