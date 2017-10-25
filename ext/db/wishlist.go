package db

import (
	"fmt"
	"strings"
)

const wishlistPut = `
INSERT INTO collections ("fk_book", "fk_user", "tags")
SELECT id, $1, '{"wishlist"}' FROM books WHERE %s ON CONFLICT DO NOTHING`

func (db *DB) WishlistPut(user int, books ...string) (int, error) {
	var args = []interface{}{user}
	var where = []string{}

	for i, isbn := range books {
		where = append(where, fmt.Sprintf("isbn = $%d", i+2))
		args = append(args, isbn)
	}
	return db.Exec(fmt.Sprintf(wishlistPut, strings.Join(where, " OR ")), args...)
}
