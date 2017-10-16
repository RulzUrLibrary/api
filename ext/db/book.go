package db

import (
	"database/sql"
	"github.com/rulzurlibrary/api/utils"
)

type Book struct {
	Id          int
	Isbn        string
	title       sql.NullString
	number      sql.NullInt64
	serie       sql.NullString
	description sql.NullString
	price       sql.NullFloat64
	owned       sql.NullBool
	authors     Authors
}

func (b *Book) ToStructs(partial bool) *utils.Book {
	book := &utils.Book{
		b.Isbn, "", b.title.String, b.description.String,
		float32(b.price.Float64), int(b.number.Int64), b.serie.String, nil,
		nil,
	}
	if !partial {
		book.Thumbnail = "/thumbs/" + b.Isbn + ".jpg"
		book.Authors = b.authors.ToStructs()
	}

	if b.owned.Valid {
		book.Owned = &b.owned.Bool
	}
	return book
}

type Books struct {
	books map[int]*Book
	order []int
}

func NewBooks() *Books {
	return &Books{make(map[int]*Book), make([]int, 0)}
}

func (b *Books) Add(book *Book) {
	if _, ok := b.books[book.Id]; !ok {
		b.order = append(b.order, book.Id)
	}
	b.books[book.Id] = book
}

func (b *Books) Get(id int) *Book {
	return b.books[id]
}

func (b *Books) First() *Book {
	if len(b.books) > 0 {
		return b.books[b.order[0]]
	}
	return nil
}

func (b *Books) ToStructs(partial bool) (books utils.Books) {
	for _, id := range b.order {
		books = append(books, b.books[id].ToStructs(partial))
	}
	return
}
