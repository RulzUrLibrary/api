package db

import (
	"database/sql"
	"fmt"
	"github.com/rulzurlibrary/api/utils"
	"strings"
)

const SelectBook = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name,
	array_agg(DISTINCT ROW(a.id, a.name)),
	array_agg(DISTINCT ROW(n.provider, n.note, n.link))
FROM books b
INNER JOIN series s ON (b.fk_serie = s.id)
LEFT OUTER JOIN notations n on (n.fk_book = b.id)
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
WHERE b.isbn = $1
GROUP BY b.id, s.name`

const SelectBookU = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name,
  c.fk_book IS NOT NULL, array_agg(DISTINCT ROW(a.id, a.name)),
	array_agg(DISTINCT ROW(w.name, w.uuid)),
  array_agg(DISTINCT ROW(n.provider, n.note, n.link))
FROM books b
INNER JOIN series s ON (b.fk_serie = s.id)
LEFT OUTER JOIN collections c ON (b.id = fk_book AND fk_user = $2)
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
LEFT OUTER JOIN wishlists_books wb on (b.id = wb.fk_book)
LEFT OUTER JOIN wishlists w on (w.id = wb.fk_wishlist)
LEFT OUTER JOIN notations n on (n.fk_book = b.id)
WHERE b.isbn = $1
GROUP BY b.id, s.name, c.fk_book`

const SelectBooks = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name,
	array_agg(DISTINCT ROW(a.id, a.name))
FROM (
	SELECT id, isbn, title, description, price, fk_serie, num
	FROM books ORDER BY id DESC LIMIT $1 OFFSET $2
) b
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
LEFT OUTER JOIN series s ON (b.fk_serie = s.id)
GROUP BY b.id, b.isbn, b.title, b.description, b.price, b,num, s.name
ORDER BY b.id DESC`

const SelectBooksU = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name,
  b.fk_book IS NOT NULL, array_agg(DISTINCT ROW(a.id, a.name))
FROM (
	SELECT id, isbn, title, description, price, fk_serie, num, fk_book
	FROM books, collections
	WHERE fk_book = id AND fk_user = $3
	ORDER BY num ASC, id DESC LIMIT $1 OFFSET $2
) b
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
LEFT OUTER JOIN series s ON (b.fk_serie = s.id)
GROUP BY b.id, b.isbn, b.title, b.description, b.price, b,num, s.name, b.fk_book
ORDER BY b.num ASC, b.id DESC`

const CountBooks = `
SELECT COUNT(*) FROM books`

const CountBooksU = `
SELECT COUNT(*) FROM books, collections WHERE fk_book = id AND fk_user = $1`

const InsertBook = `
INSERT INTO books (isbn, title, description, price, num, fk_serie)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`

const InsertAuthor = `
WITH s AS (
	SELECT id FROM authors WHERE lower(name) = lower($1)
), i AS (
	INSERT INTO authors ("name") SELECT $1 WHERE NOT EXISTS (SELECT 1 FROM s) RETURNING id
)
SELECT id FROM i UNION ALL SELECT id FROM s`

const InsertBookAuthor = `
INSERT INTO book_authors (fk_book, fk_author)
VALUES ($1, $2)`

const InsertBookNotation = `
INSERT INTO notations (fk_book, provider, note, link) VALUES ($1, $2, $3, $4)`

const InsertSerie = `
WITH s AS (
	SELECT id FROM series WHERE lower(name) = lower($1)
), i AS (
	INSERT INTO series ("name") SELECT $1 WHERE NOT EXISTS (SELECT 1 FROM s) RETURNING id
)
SELECT id FROM i UNION ALL SELECT id FROM s`

const inCollection = `
SELECT EXISTS(
	SELECT 1 FROM books b, collections c WHERE c.fk_user = $1 AND c.fk_book = $2
)`

const SelectCollection = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name, a.name
FROM (
	SELECT id, isbn, title, description, price, num, fk_serie
	FROM books, collections
	WHERE fk_user = $3 AND fk_book = id LIMIT $1 OFFSET $2
) b
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
LEFT OUTER JOIN series s ON (b.fk_serie = s.id)`

const CountCollection = `
SELECT COUNT(id) FROM books, collections WHERE fk_user = $1 AND fk_book = id`

const SearchBooks = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name,
	array_agg(DISTINCT ROW(a.id, a.name))
FROM (
  SELECT id, isbn, num, title, description, price, fk_serie,
    ts_rank_cd(tsv, query) AS rank
  FROM books, plainto_tsquery('french', $3) query
  WHERE query @@ tsv
  ORDER BY num NULLS FIRST, rank DESC
  LIMIT $1 OFFSET $2
) b
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
LEFT OUTER JOIN series s ON (b.fk_serie = s.id)
GROUP BY b.id, b.isbn, b.title, b.description, b.price, b.num, s.name
ORDER BY num NULLS FIRST`

const InsertCollection = `
INSERT INTO collections ("fk_book", "fk_user")
SELECT id, $1 FROM books WHERE %s ON CONFLICT DO NOTHING`

const DeleteCollection = `
DELETE FROM collections USING books b
WHERE fk_user = $1 AND b.id = fk_book AND (%s)`

func (db *DB) BookGet(id string) (book Book, err error) {
	err = db.QueryRow(SelectBook, id).Scan(
		&book.id, &book.isbn, &book.title, &book.description, &book.price,
		&book.number, &book.serie, &book.authors, &book.notations,
	)
	if err == sql.ErrNoRows {
		err = utils.ErrNotFound
	}
	return

}

func (db *DB) BookGetU(id string, user int) (book Book, err error) {
	err = db.QueryRow(SelectBookU, id, user).Scan(
		&book.id, &book.isbn, &book.title, &book.description, &book.price,
		&book.number, &book.serie, &book.owned, &book.authors, &book.wishlists,
		&book.notations,
	)
	if err == sql.ErrNoRows {
		err = utils.ErrNotFound
	}
	return
}

func (db *DB) BookSave(book *utils.Book) error {

	args := list{
		book.Isbn, toInterfaceS(book.Title), book.Description, book.Price,
		toInterfaceI(book.Number), nil,
	}

	return db.Transaction(func(tx *sql.Tx) (err error) {
		var idBook, idAuthor int
		var insert = func(query string, args ...interface{}) (id int, err error) {
			err = tx.QueryRow(query, args...).Scan(&id)
			return
		}
		args[5], err = insert(InsertSerie, toInterfaceS(book.Serie))
		if err != nil {
			return
		}

		if idBook, err = insert(InsertBook, args...); err != nil {
			return
		}

		for _, author := range *book.Authors {
			if idAuthor, err = insert(InsertAuthor, author.Name); err != nil {
				return
			}

			if _, err = tx.Exec(InsertBookAuthor, idBook, idAuthor); err != nil {
				return
			}
		}
		for _, notation := range *book.Notations {
			a := list{idBook, notation.Provider, notation.Note, notation.Link}
			if _, err = tx.Exec(InsertBookNotation, a...); err != nil {
				return
			}
		}
		return
	})
}

func (db *DB) BookList(limit, offset int) (books Books, count int64, err error) {
	return db.queryList(queryList{
		query{SelectBooks, list{limit, offset}, func(b *Book) list {
			return list{&b.id, &b.isbn, &b.title, &b.description, &b.price,
				&b.number, &b.serie, &b.authors}
		}}, CountBooks, list{},
	})
}

func (db *DB) BookListU(limit, offset, user int) (books Books, count int64, err error) {
	return db.queryList(queryList{
		query{SelectBooksU, list{limit, offset, user}, func(b *Book) list {
			return list{&b.id, &b.isbn, &b.title, &b.description, &b.price,
				&b.number, &b.serie, &b.owned, &b.authors}
		}}, CountBooksU, list{user},
	})
}

func (db *DB) BookSearch(pattern string, limit, offset int) (books Books, err error) {
	return db.query(query{SearchBooks, list{limit, offset, pattern}, func(b *Book) list {
		return list{&b.id, &b.isbn, &b.title, &b.description, &b.price,
			&b.number, &b.serie, &b.authors}
	}})
}

func (db *DB) BookDelete(user int, books ...string) (int64, error) {
	var args = list{user}
	var where = []string{}

	for i, isbn := range books {
		where = append(where, fmt.Sprintf("b.isbn = $%d", i+2))
		args = append(args, isbn)
	}

	return db.Exec(fmt.Sprintf(DeleteCollection, strings.Join(where, " OR ")), args...)
}

func (db *DB) BookPut(user int, books ...string) (int64, error) {
	var args = list{user}
	var where = []string{}

	for i, isbn := range books {
		where = append(where, fmt.Sprintf("isbn = $%d", i+2))
		args = append(args, isbn)
	}
	return db.Exec(fmt.Sprintf(InsertCollection, strings.Join(where, " OR ")), args...)
}
