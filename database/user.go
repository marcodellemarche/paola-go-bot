package database

import (
	"log"

	_ "embed"

	_ "github.com/lib/pq"
)

//go:embed user/create.sql
var sqlUserCreateTable string

//go:embed user/drop.sql
var sqlUserDropTable string

//go:embed user/insert.sql
var sqlUserInsert string

//go:embed user/delete_by_id.sql
var sqlUserDeleteById string

//go:embed user/find_by_id.sql
var sqlUserFindById string

type User struct {
	Id      int64
	Name    string
}

func UserNew(Id int64, Name string) User {
	return User {
		Id,
		Name,
	}
}

func UserCreateTable() {
	_, err := db.Exec(sqlUserCreateTable)
	CheckError(err)
}

func UserDropTable() {
	_, err := db.Exec(sqlUserDropTable)
	CheckError(err)
}

func UserInsert(id int64, name string) bool {
	_, err := db.Exec(sqlUserInsert, id, name)
	if err != nil {
		log.Printf("Error inserting user into database, user insertion failed: %s", err.Error())
		return false
	}

	return true
}

func UserFindById(id int64) (User, bool) {
	var user User

	rows, err := db.Query(sqlUserFindById, id)
	if err != nil {
		log.Printf("Error fetching user from database, quary failed: %s", err.Error())
		return user, false
	}

	if !rows.Next() {
		log.Printf("No user %d found", id)
		return user, true
	}

	defer rows.Close()

	rows.Scan(&user.Id, &user.Name)
	if err != nil {
		log.Printf("Error fetching user from database, scan failed: %s", err.Error())
		return user, false
	}

	return user, true
}

func UserDeleteById(id int64) {
	_, err := db.Exec(sqlUserDeleteById, id)
	CheckError(err)
}
