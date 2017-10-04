package db

import (
	"database/sql"
	"encoding/json"
	"github.com/paul-bismuth/library/utils"
)

type Book struct {
	Id          int
	Isbn        string
	title       sql.NullString
	number      sql.NullInt64
	serie       sql.NullString
	description sql.NullString
	price       sql.NullFloat64
	owned       bool
	authors     Authors
}

func (b *Book) Title() (title string) {
	return b.title.String
}

func (b *Book) TitleDisplay() string {
	return utils.TitleDisplay(b.Title(), b.Serie(), b.Number())
}

func (b *Book) Description() string {
	return b.description.String
}

func (b *Book) Thumbnail() string {
	return "/thumbs/" + b.Isbn + ".jpg"
}

func (b *Book) Price() float32 {
	return float32(b.price.Float64)
}

func (b *Book) Number() int {
	return int(b.number.Int64)
}

func (b *Book) Serie() string {
	return b.serie.String
}

func (b *Book) Authors() utils.Authors {
	return b.authors.ToStructs()
}

func (b *Book) ToBook() *utils.Book {
	return &utils.Book{
		Isbn:        b.Isbn,
		Title:       b.Title(),
		Description: b.Description(),
		Thumbnail:   b.Thumbnail(),
		Price:       b.Price(),
		Authors:     b.Authors(),
		Number:      b.Number(),
		Serie:       b.Serie(),
	}
}

func (b *Book) ToBookScoped() *utils.BookScoped {
	return &utils.BookScoped{*b.ToBook(), b.owned}
}

func (b *Book) ToVolume() *utils.Volume {
	return &utils.Volume{
		Isbn:   b.Isbn,
		Number: b.Number(),
		Title:  b.Title(),
	}
}

func (b *Book) ToVolumeScoped() *utils.VolumeScoped {
	return &utils.VolumeScoped{*b.ToVolume(), b.owned}
}

func (b *Book) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.ToBook())
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

func (b *Books) ToBooks() (books []*utils.Book) {
	for _, id := range b.order {
		books = append(books, b.books[id].ToBook())
	}
	return
}

func (b *Books) ToBooksScoped() (books []*utils.BookScoped) {
	for _, id := range b.order {
		books = append(books, b.books[id].ToBookScoped())
	}
	return
}

func (b *Books) ToVolumes() (volumes utils.Volumes) {
	for _, id := range b.order {
		volumes = append(volumes, b.books[id].ToVolume())
	}
	return
}

func (b *Books) ToVolumesScoped() (volumes utils.VolumesScoped) {
	for _, id := range b.order {
		volumes = append(volumes, b.books[id].ToVolumeScoped())
	}
	return
}
