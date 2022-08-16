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

//go:embed user/find_all.sql
var sqlUserFindAll string

type User struct {
	Id       int64
	Name     string
}

func UserCreateTable() bool {
	_, err := db.Exec(sqlUserCreateTable)
	if err != nil {
		log.Printf("Error creating user table: %s", err.Error())
		return false
	}

	return true
}

func UserDropTable() bool {
	_, err := db.Exec(sqlUserDropTable)
	if err != nil {
		log.Printf("Error dropping user table: %s", err.Error())
		return false
	}

	return true
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
		log.Printf("Error fetching user from database, query failed: %s", err.Error())
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

func UserFindAll() ([]User, bool) {
	rows, err := db.Query(sqlUserFindAll)
	if err != nil {
		log.Printf("Error finding users from database, query failed: %s", err.Error())
		return nil, false
	}

	users := make([]User, 0)
	for rows.Next() {
		var user User

		err := rows.Scan(&user.Id, &user.Name)
		if err != nil {
			log.Printf("Error finding users from database, scan failed: %s", err.Error())
			return nil, false
		}

		users = append(users, user)
	}

	return users, true
}

func UserDeleteById(id int64) bool {
	_, err := db.Exec(sqlUserDeleteById, id)
	if err != nil {
		log.Printf("Error deleting user from database, deletion failed: %s", err.Error())
		return false
	}

	return true
}
