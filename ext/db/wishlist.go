package db

import (
	"fmt"
	"github.com/rulzurlibrary/api/utils"
	"strings"
)

const wishlistPut = `
INSERT INTO collections ("fk_book", "fk_user", "tags")
SELECT id, $1, '{"wishlist"}' FROM books WHERE %s ON CONFLICT DO NOTHING`

const countWishList = `
SELECT COUNT(*) FROM collections WHERE fk_user = $1 AND 'wishlist'=ANY(tags)`

const wishList = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name, a.id,
	a.name, tags
FROM (
	SELECT id, isbn, title, description, price, fk_serie, num, tags
	FROM books, collections
	WHERE fk_book = id AND fk_user = $3 AND 'wishlist'=ANY(tags)
	ORDER BY num ASC, id DESC LIMIT $1 OFFSET $2
) b
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
LEFT OUTER JOIN series s ON (b.fk_serie = s.id)
ORDER BY b.num ASC, b.id DESC`

func (db *DB) WishlistPut(user int, books ...string) (int, error) {
	var args = []interface{}{user}
	var where = []string{}

	for i, isbn := range books {
		where = append(where, fmt.Sprintf("isbn = $%d", i+2))
		args = append(args, isbn)
	}
	return db.Exec(fmt.Sprintf(wishlistPut, strings.Join(where, " OR ")), args...)
}

func (db *DB) WishList(limit, offset, user int) ([]*utils.Book, int, error) {
	return db.bookList(queryBookList{
		queryBook{wishList, []interface{}{limit, offset, user},
			func(b *Book, a *Author) []interface{} {
				return []interface{}{
					&b.Id, &b.Isbn, &b.title, &b.description, &b.price, &b.number,
					&b.serie, &a.id, &a.name, &b.tags,
				}
			},
		},
		countWishList, []interface{}{user},
	})
}
