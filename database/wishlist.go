package database

import (
	"log"

	_ "embed"

	"database/sql"

	_ "github.com/lib/pq"
)

//go:embed wishlist/create.sql
var sqlWishlistCreateTable string

//go:embed wishlist/drop.sql
var sqlWishlistDropTable string

//go:embed wishlist/insert.sql
var sqlWishlistInsert string

//go:embed wishlist/delete_by_name.sql
var sqlWishlistDeleteByName string

//go:embed wishlist/find_by_user.sql
var sqlWishlistFindByUser string

//go:embed wishlist/set_buyer.sql
var sqlWishlistSetBuyer string

type Wishlist struct {
	UserId  int64
	Name    string
	Link    string
	BuyerId int64
}

func WishlistCreateTable() bool {
	_, err := db.Exec(sqlWishlistCreateTable)
	if err != nil {
		log.Printf("Error creating wishlist table: %s", err.Error())
		return false
	}

	return true
}

func WishlistDropTable() bool {
	_, err := db.Exec(sqlWishlistDropTable)
	if err != nil {
		log.Printf("Error dropping wishlist table: %s", err.Error())
		return false
	}

	return true
}

func WishlistInsert(userId int64, name string, link string, buyerId int64) bool {
	nullableBuyerId := sql.NullInt64{Int64: buyerId, Valid: buyerId > 0}
	nullableLink := sql.NullString{String: link, Valid: link != ""}

	_, err := db.Exec(sqlWishlistInsert, userId, name, nullableLink, nullableBuyerId)
	if err != nil {
		log.Printf("Error inserting wishlist into database, wishlist insertion failed: %s", err.Error())
		return false
	}

	return true
}

func WishlistFindByUser(userId int64) ([]Wishlist, bool) {
	rows, err := db.Query(sqlWishlistFindByUser, userId)
	if err != nil {
		log.Printf("Error finding wishlists from database, query failed: %s", err.Error())
		return nil, false
	}

	return wishlistFind(rows)
}

func wishlistFind(rows *sql.Rows) ([]Wishlist, bool) {
	wishlists := make([]Wishlist, 0)
	for rows.Next() {
		var wishlist Wishlist

		var nullableBuyerId sql.NullInt64
		var nullableLink sql.NullString

		err := rows.Scan(&wishlist.UserId, &wishlist.Name, &nullableLink, &nullableBuyerId)
		if err != nil {
			log.Printf("Error finding wishlists from database, scan failed: %s", err.Error())
			return nil, false
		}

		if nullableLink.Valid {
			wishlist.Link = nullableLink.String
		}

		if nullableBuyerId.Valid {
			wishlist.BuyerId = nullableBuyerId.Int64
		}

		wishlists = append(wishlists, wishlist)
	}

	return wishlists, true
}

func WishlistDeleteByName(name int64, userId int64) bool {
	_, err := db.Exec(sqlWishlistDeleteByName, name, userId)
	if err != nil {
		log.Printf("Error deleting wishlist from database, wishlist deletion failed: %s", err.Error())
		return false
	}

	return true
}

func WishlistSetBuyer(name string, userId int64, buyerId int64) bool {
	_, err := db.Exec(sqlWishlistSetBuyer, name, userId, buyerId)
	if err != nil {
		log.Printf("Error setting buyer for wishlist on database, wishlist update failed: %s", err.Error())
		return false
	}

	return true
}
