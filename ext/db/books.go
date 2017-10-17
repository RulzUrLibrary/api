package db

import (
	"fmt"
	"github.com/rulzurlibrary/api/utils"
	"strings"
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

const SelectBooksU = `
SELECT b.id, b.isbn, b.title, b.description, b.price, b.num, s.name, a.id, a.name
FROM (
	SELECT id, isbn, title, description, price, fk_serie, num
	FROM books, collections
	WHERE fk_book = id AND fk_user = $3
	ORDER BY num ASC, id DESC LIMIT $1 OFFSET $2
) b
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
LEFT OUTER JOIN series s ON (b.fk_serie = s.id)
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

const InsertCollection = `
INSERT INTO collections ("fk_book", "fk_user")
SELECT id, $1 FROM books WHERE %s ON CONFLICT DO NOTHING`

const DeleteCollection = `
DELETE FROM collections USING books b
WHERE fk_user = $1 AND b.id = fk_book AND (%s)`

type queryBook struct {
	query   string
	args    []interface{}
	getArgs func(*Book, *Author) []interface{}
}

type queryBookList struct {
	queryBook
	queryList     string
	queryListArgs []interface{}
}

func dedupBooks(db *DB, qb queryBook) (*Books, error) {
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

func (db *DB) BookGet(id string) (*Book, error) {
	return db.bookGet(queryBook{
		SelectBook, []interface{}{id},
		func(b *Book, a *Author) []interface{} {
			return []interface{}{
				&b.Id, &b.Isbn, &b.title, &b.description, &b.price, &b.number,
				&b.serie, &a.id, &a.name,
			}
		},
	})
}

func (db *DB) BookGetU(id string, user int) (*Book, error) {
	return db.bookGet(queryBook{
		SelectBookU, []interface{}{id, user},
		func(b *Book, a *Author) []interface{} {
			return []interface{}{
				&b.Id, &b.Isbn, &b.title, &b.description, &b.price, &b.number,
				&b.serie, &a.id, &a.name, &b.owned,
			}
		},
	})
}

func (db *DB) bookGet(qb queryBook) (*Book, error) {
	books, err := dedupBooks(db, qb)
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
		var idBook, idAuthor int
		args[5], err = tx.Insert(InsertSerie, toInterfaceS(book.Serie))
		if err != nil {
			return
		}

		if idBook, err = tx.Insert(InsertBook, args...); err != nil {
			return
		}

		for _, author := range *book.Authors {
			if idAuthor, err = tx.Insert(InsertAuthor, author.Name); err != nil {
				return
			}

			if _, err = tx.Exec(InsertBookAuthor, idBook, idAuthor); err != nil {
				return
			}
		}
		return
	})
}

func (db *DB) BookList(limit, offset int) ([]*utils.Book, int, error) {
	return db.bookList(queryBookList{
		queryBook{
			SelectBooks, []interface{}{limit, offset},
			func(b *Book, a *Author) []interface{} {
				return []interface{}{
					&b.Id, &b.Isbn, &b.title, &b.description, &b.price, &b.number,
					&b.serie, &a.id, &a.name,
				}
			},
		},
		CountBooks, []interface{}{},
	})
}

func (db *DB) BookListU(limit, offset, user int) ([]*utils.Book, int, error) {
	return db.bookList(queryBookList{
		queryBook{
			SelectBooksU, []interface{}{limit, offset, user},
			func(b *Book, a *Author) []interface{} {
				return []interface{}{
					&b.Id, &b.Isbn, &b.title, &b.description, &b.price, &b.number,
					&b.serie, &a.id, &a.name,
				}
			},
		},
		CountBooksU, []interface{}{user},
	})
}

func (db *DB) BookSearch(pattern string, limit, offset int) ([]*utils.Book, error) {
	books, err := dedupBooks(db, queryBook{
		SearchBooks, []interface{}{limit, offset, pattern},
		func(b *Book, a *Author) []interface{} {
			return []interface{}{
				&b.Id, &b.Isbn, &b.title, &b.description, &b.price, &b.number,
				&b.serie, &a.id, &a.name,
			}
		},
	})
	if err != nil {
		return nil, err
	}
	return books.ToStructs(false), nil
}

func (db *DB) bookList(qbl queryBookList) ([]*utils.Book, int, error) {
	count, err := db.Count(qbl.queryList, qbl.queryListArgs...)
	if err != nil {
		return nil, 0, err
	}
	books, err := dedupBooks(db, qbl.queryBook)
	if err != nil {
		return nil, 0, err
	}
	return books.ToStructs(false), count, nil
}

func (db *DB) BookDelete(user int, books ...string) (int, error) {
	var args = []interface{}{user}
	var where = []string{}

	for i, isbn := range books {
		where = append(where, fmt.Sprintf("b.isbn = $%d", i+2))
		args = append(args, isbn)
	}

	return db.Exec(fmt.Sprintf(DeleteCollection, strings.Join(where, " OR ")), args...)
}

func (db *DB) BookPut(user int, books ...string) (int, error) {
	var args = []interface{}{user}
	var where = []string{}

	for i, isbn := range books {
		where = append(where, fmt.Sprintf("isbn = $%d", i+2))
		args = append(args, isbn)
	}
	return db.Exec(fmt.Sprintf(InsertCollection, strings.Join(where, " OR ")), args...)
}
