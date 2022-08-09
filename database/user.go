package database

import (
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

func UserCreateTable() {
	_, err := db.Exec(sqlUserCreateTable)
	CheckError(err)
}

func UserDropTable() {
	_, err := db.Exec(sqlUserDropTable)
	CheckError(err)
}

func UserInsert(id int64, name string) {
	_, err := db.Exec(sqlUserInsert, id, name)
	CheckError(err)
}

func UserDeleteById(id int64) {
	_, err := db.Exec(sqlUserDeleteById, id)
	CheckError(err)
}
