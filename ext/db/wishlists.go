package db

import (
	"fmt"
	"github.com/RulzUrLibrary/api/utils"
	"strings"
)

const wishlistPut = `
INSERT INTO wishlists_books ("fk_book", "fk_wishlist")
SELECT b.id, w.id
FROM wishlists w, books b, users u
WHERE w.fk_user = $1 AND b.isbn = $2 AND (%s)
ON CONFLICT DO NOTHING`

const countWishlists = `
SELECT COUNT(id) FROM wishlists WHERE fk_user = $1`

const countWishlist = `
SELECT COUNT(fk_book) FROM wishlists, wishlists_books WHERE fk_wishlist = id AND uuid = $1`

const wishlists = `
SELECT w.id, w.name, w.uuid, b.id, b.isbn, b.title, b.price, b.num, s.name,
	array_agg(DISTINCT ROW(a.id, a.name))
FROM (
	SELECT id, name, uuid FROM wishlists WHERE fk_user = $3
	ORDER BY id DESC LIMIT $1 OFFSET $2
) w
FULL JOIN wishlists_books wb ON (w.id = wb.fk_wishlist)
LEFT OUTER JOIN books b ON (b.id = wb.fk_book)
LEFT OUTER JOIN series s ON (s.id = b.fk_serie)
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
GROUP BY w.id, w.name, w.uuid, s.name, b.id`

const wishlist = `
SELECT w.id, w.name, w.uuid, u.email, b.id, b.isbn, b.title, b.price, b.num,
	s.name, array_agg(DISTINCT ROW(a.id, a.name))
FROM wishlists w
LEFT OUTER JOIN users u ON (w.fk_user = u.id)
FULL JOIN wishlists_books wb ON (w.id = wb.fk_wishlist)
LEFT OUTER JOIN books b ON (b.id = wb.fk_book)
LEFT OUTER JOIN series s ON (s.id = b.fk_serie)
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
WHERE w.uuid = $3 GROUP BY w.id, w.name, w.uuid, u.email, s.name, b.id LIMIT $1 OFFSET $2`

const wishlistsN = `
SELECT w.id, w.name, w.uuid FROM wishlists w WHERE w.fk_user = $1`

const wishInsert = `
INSERT INTO wishlists (uuid, name, fk_user) VALUES (gen_random_uuid(), $1, $2)`

const wishlistClean = `
DELETE FROM wishlists_books USING wishlists w, books b
WHERE fk_wishlist = w.id AND fk_user = $1 AND fk_book = b.id AND b.isbn = $2`

const wishlistDelete = `DELETE FROM wishlists WHERE fk_user = $1 AND (%s)`

func (db *DB) WishlistPut(user int, book string, wishlists ...string) (int64, error) {
	var args = list{user, book}
	var where = []string{}

	for i, wishlist := range wishlists {
		where = append(where, fmt.Sprintf("w.uuid = $%d", i+3))
		args = append(args, wishlist)
	}
	return db.Exec(fmt.Sprintf(wishlistPut, strings.Join(where, " OR ")), args...)
}

func (db *DB) Wishlist(limit, offset int, uuid string) (w Wishlists, c int64, e error) {
	c, e = db.queryList(wishlist, list{limit, offset, uuid}, w.InsertWishlist(func(w *Wishlist) list {
		return list{&w.id, &w.name, &w.uuid, &w.user, &w.book.id, &w.book.isbn, &w.book.title,
			&w.book.price, &w.book.number, &w.book.serie, &w.book.authors,
		}
	}), countWishlist, list{uuid})
	return
}

func (db *DB) Wishlists(limit, offset, user int) (w Wishlists, c int64, e error) {
	c, e = db.queryList(wishlists, list{limit, offset, user}, w.InsertWishlist(func(w *Wishlist) list {
		return list{&w.id, &w.name, &w.uuid, &w.book.id, &w.book.isbn, &w.book.title,
			&w.book.price, &w.book.number, &w.book.serie, &w.book.authors}
	}), countWishlists, list{user})
	return
}

func (db *DB) WishlistsN(user int) (w Wishlists, e error) {
	e = db.query(wishlistsN, list{user}, w.InsertWishlist(func(w *Wishlist) list {
		return list{&w.id, &w.name, &w.uuid}
	}))
	return
}

func (db *DB) WishPost(name string, user int) error {
	count, err := db.Exec(wishInsert, name, user)
	if err != nil {
		return err
	}
	if count == 0 {
		return utils.ErrExists
	}
	return nil
}

func (db *DB) WishlistDelete(user int, uuids ...string) (int64, error) {
	var args = list{user}
	var where = []string{}

	for i, uuid := range uuids {
		where = append(where, fmt.Sprintf("uuid = $%d", i+2))
		args = append(args, uuid)
	}
	return db.Exec(fmt.Sprintf(wishlistDelete, strings.Join(where, " OR ")), args...)
}

func (db *DB) WishlistUpdate(user int, book string, uuids ...string) error {
	return db.Transaction(func(tx *Tx) error {
		if _, err := tx.Exec(wishlistClean, user, book); err != nil {
			return err
		}
		if len(uuids) == 0 {
			return nil
		}
		args := list{user, book}
		where := []string{}

		for i, uuid := range uuids {
			where = append(where, fmt.Sprintf("w.uuid = $%d", i+3))
			args = append(args, uuid)
		}
		if _, err := tx.Exec(fmt.Sprintf(wishlistPut, strings.Join(where, " OR ")), args...); err != nil {
			return err
		}
		return nil
	})
}
