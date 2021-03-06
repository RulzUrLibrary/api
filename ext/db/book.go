package db

import (
	"database/sql"
	"github.com/RulzUrLibrary/api/utils"
)

type Book struct {
	id          sql.NullInt64
	isbn        sql.NullString
	title       sql.NullString
	number      sql.NullInt64
	serie       sql.NullString
	serie_id    sql.NullInt64
	description sql.NullString
	price       sql.NullFloat64
	owned       sql.NullBool
	authors     Authors
	wishlists   Wishlists
	notations   Notations
}

func (b Book) ToStructs(partial bool) (book utils.Book) {
	book.Isbn = b.isbn.String
	book.Title = b.title.String
	book.Price = float32(b.price.Float64)
	book.Serie = b.serie.String
	book.Number = int(b.number.Int64)
	book.Wishlists = b.wishlists.ToWishlists(partial)

	if !partial {
		book.Description = b.description.String
		book.Thumbnail = "/thumbs/" + b.isbn.String + ".jpg"
		book.Authors = b.authors.Authors
		book.Notations = b.notations.Notations
	}

	if b.owned.Valid {
		book.Owned = &b.owned.Bool
	}

	return
}

func (b Book) ToSerie(partial bool) (serie utils.Serie) {
	serie.Description = b.description.String
	serie.Authors = b.authors.Authors

	if b.number.Valid {
		serie.Id = b.serie_id.Int64
		serie.Name = b.serie.String

		b.serie.String = ""
		serie.Volumes = &utils.Books{b.ToStructs(partial)}
	} else {
		serie.Isbn = b.isbn.String
		serie.Title = b.title.String
		if b.owned.Valid {
			serie.Owned = &b.owned.Bool
		}
	}
	return
}

type Books struct {
	Books []*Book
}

func (b *Books) ToStructs(partial bool) utils.Books {
	// https://apoorvam.github.io/golang/json/marshal/slice/empty/null/2017/01/19/golang-json-marshalling.html
	books := make(utils.Books, len(b.Books))
	for i, book := range b.Books {
		books[i] = book.ToStructs(partial)
	}
	return books
}

func (b *Books) InsertBook(fn func(b *Book) list) func() list {
	return func() list {
		var book Book
		b.Books = append(b.Books, &book)
		return fn(&book)
	}
}

func (b *Books) ToSeries(partial bool) (series utils.Series) {
	var last utils.Serie

	for _, book := range b.Books {
		if book.serie_id.Int64 == last.Id {
			book.serie.String = "" // we are dumping series, no need to keep this
			*last.Volumes = append(*last.Volumes, book.ToStructs(partial))
		} else {
			last = book.ToSerie(partial)
			series = append(series, last)
		}
	}
	return series
}
