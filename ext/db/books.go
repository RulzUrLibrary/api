package db

import (
	//"database/sql"
	"github.com/ixday/echo-hello/utils"
)

const SelectBook = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name, a.id, a.name
FROM books b
INNER JOIN series s ON (b.fk_serie = s.id)
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
WHERE b.isbn = $1`

const SelectBookU = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name, a.id, a.name,
	EXISTS(SELECT true FROM collections WHERE fk_book = b.id AND fk_user = $2)
FROM books b
INNER JOIN series s ON (b.fk_serie = s.id)
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
WHERE b.isbn = $1`

const SelectBooks = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name, a.id, a.name
FROM (
	SELECT id, isbn, title, description, price, fk_serie, num
	FROM books ORDER BY id DESC LIMIT $1 OFFSET $2
) b
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
LEFT OUTER JOIN series s ON (b.fk_serie = s.id)
ORDER BY b.id DESC`

const CountBooks = `
SELECT COUNT(*) FROM books`

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

const AddBook = `
INSERT INTO collections SELECT b.id, $2 FROM books b WHERE b.isbn = $1
ON CONFLICT (fk_book, fk_user) DO UPDATE SET fk_user = $2`

const SearchBooks = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name, a.id, a.name
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
ORDER BY num NULLS FIRST`

type queryBook struct {
	query   string
	args    []interface{}
	getArgs func(*Book, *Author) []interface{}
}

func dedup(db *DB, qb queryBook) (*Books, error) {
	books := NewBooks()
	rows, err := db.Query(qb.query, qb.args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var book Book
		var author Author

		if err := rows.Scan(qb.getArgs(&book, &author)...); err != nil {
			return nil, err
		}

		if b := books.Get(book.Id); b != nil {
			b.authors = append(b.authors, author)
		} else {
			book.authors = Authors{author}
			books.Add(&book)
		}
	}
	return books, nil
}

func newBookQuery(id string, user int) (qb queryBook) {
	if user == 0 {
		qb = queryBook{
			SelectBook, []interface{}{id},
			func(b *Book, a *Author) []interface{} {
				return []interface{}{
					&b.Id, &b.Isbn, &b.title, &b.description, &b.price, &b.number,
					&b.serie, &a.id, &a.name,
				}
			},
		}
	} else {
		qb = queryBook{
			SelectBookU, []interface{}{id, user},
			func(b *Book, a *Author) []interface{} {
				return []interface{}{
					&b.Id, &b.Isbn, &b.title, &b.description, &b.price, &b.number,
					&b.serie, &a.id, &a.name, &b.owned,
				}
			},
		}
	}
	return
}

func (db *DB) BookGet(id string, user int) (*Book, error) {
	books, err := dedup(db, newBookQuery(id, user))
	if err != nil {
		return nil, err
	}
	book := books.First()
	if book == nil {
		return nil, utils.ErrNotFound
	}

	return book, nil
}

func (db *DB) BookSave(book *utils.Book) error {
	args := []interface{}{
		book.Isbn, toInterfaceS(book.Title), book.Description, book.Price,
		toInterfaceI(book.Number), nil,
	}

	return db.Transaction(func(tx *Tx) (err error) {
		args[5], err = tx.Insert(InsertSerie, toInterfaceS(book.Serie))
		if err != nil {
			return
		}

		if book.Id, err = tx.Insert(InsertBook, args...); err != nil {
			return
		}

		for _, author := range book.Authors {
			id, err := tx.Insert(InsertAuthor, author.Name)
			if err != nil {
				return err
			}
			_, err = tx.Exec(InsertBookAuthor, book.Id, id)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

//func BookList(query string, args ...interface{}) ([]*utils.Book, error) {
//	books, err := dedup(query, args...)
//	if err != nil {
//		return nil, err
//	}
//	return books.Gets(), nil
//}
//
//func BookAdd(user int, isbn string) error {
//	if len(isbn) == 0 {
//		return utils.ErrInvalidIsbn
//	}
//	result, err := Client.Exec(AddBook, isbn, user)
//	if err != nil {
//		return nil
//	}
//	count, err := result.RowsAffected()
//	if err != nil {
//		return err
//	}
//	if int(count) == 0 {
//		return utils.ErrNoProduct
//	}
//	return err
//}
//
//func BookSave(book *utils.Book) error {
//	var args []interface{}
//	var id int
//
//	return Transaction(func(tx *Tx) (err error) {
//		args = book.ToArgs()
//
//		if book.Serie == "" {
//			args[5], err = tx.Insert(InsertSerie, nil)
//		} else {
//			args[5], err = tx.Insert(InsertSerie, book.Serie)
//		}
//		if err != nil {
//			return
//		}
//
//		if book.Id, err = tx.Insert(InsertBook, args...); err != nil {
//			return
//		}
//
//		for _, author := range book.Authors {
//			if id, err = tx.Insert(InsertAuthor, author.Name); err != nil {
//				return
//			}
//			err = tx.QueryRow(InsertBookAuthor, book.Id, id).Scan()
//			if err != nil && err != sql.ErrNoRows {
//				return
//			}
//		}
//
//		return nil
//	})
//}
//
//func InCollection(book int, u *utils.User) (ok bool, err error) {
//	if u == nil {
//		return
//	}
//	err = Client.QueryRow(inCollection, u.Id, book).Scan(&ok)
//	return
//}
