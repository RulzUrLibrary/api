package db

import (
	"github.com/paul-bismuth/library/utils"
)

const CountSeries = `
SELECT COUNT(DISTINCT(id))
FROM (
	SELECT s.id FROM books b, collections, series s
	WHERE fk_user = $1 AND fk_book = b.id AND b.fk_serie = s.id
) _`

const SelectSeries = `
SELECT title, num, isbn, s.id, s.name, a.id, a.name, EXISTS(
  SELECT true FROM collections WHERE fk_book = b.id AND fk_user = $3
), (SELECT b.description WHERE b.num IS NULL OR b.num = 1)
FROM (
    SELECT s.id FROM books b, collections, series s
    WHERE fk_user = $3 AND fk_book = b.id AND fk_serie = s.id
    GROUP BY s.id ORDER BY s.name DESC NULLS LAST LIMIT $1 OFFSET $2
) r, books b
INNER JOIN series s ON (fk_serie = s.id)
LEFT OUTER JOIN book_authors ba ON (b.id = fk_book)
LEFT JOIN authors a ON (fk_author = a.id)
WHERE s.id = r.id ORDER BY s.name, num`

const SelectSerie = `
SELECT b.id, title, num, description, isbn, a.id, a.name, s.name
FROM books b
INNER JOIN series s ON (fk_serie = s.id)
LEFT OUTER JOIN book_authors ba ON (b.id = fk_book)
LEFT JOIN authors a ON (fk_author = a.id)
WHERE s.id = $1 ORDER BY num`

const SelectSerieU = `
SELECT b.id, title, num, description, isbn, a.id, a.name, s.name, EXISTS(
  SELECT true FROM collections WHERE fk_book = b.id AND fk_user = $2
)
FROM books b
INNER JOIN series s ON (fk_serie = s.id)
LEFT OUTER JOIN book_authors ba ON (b.id = fk_book)
LEFT JOIN authors a ON (fk_author = a.id)
WHERE s.id = $1 ORDER BY num`

type querySerie struct {
	query   string
	args    []interface{}
	getArgs func(*Serie, *Book, *Author) []interface{}
}

func (db *DB) SerieList(user, limit, offset int) ([]*utils.Serie, error) {
	series := NewSeries()

	rows, err := db.Query(SelectSeries, limit, offset, user)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var serie Serie
		var author Author
		var volume Volume

		err = rows.Scan(
			&serie.title, &volume.number, &volume.isbn, &serie.id, &serie.name,
			&author.id, &author.name, &volume.owned, &serie.description,
		)
		if err != nil {
			return nil, err
		}

		if s := series.Get(serie.id); s != nil {
			s.authors = append(s.authors, author)
			s.volumes = append(s.volumes, volume)
		} else {
			serie.authors = Authors{author}
			serie.volumes = Volumes{volume}
			series.Add(&serie)
		}
	}
	return series.Gets(), nil
}

func dedupSerie(db *DB, qs querySerie) (*utils.Serie, *Books, error) {
	books := NewBooks()
	serie := Serie{}

	rows, err := db.Query(qs.query, qs.args...)

	if err != nil {
		return nil, nil, err
	}

	for rows.Next() {
		var book Book
		var author Author

		if err := rows.Scan(qs.getArgs(&serie, &book, &author)...); err != nil {
			return nil, nil, err
		}

		serie.authors = append(serie.authors, author)
		if b := books.Get(book.Id); b != nil {
			b.authors = append(b.authors, author)
		} else {
			book.serie = serie.name
			book.authors = Authors{author}
			books.Add(&book)
		}
	}
	if serie.name.String == "" {
		return nil, nil, utils.ErrNotFound
	}
	return serie.ToStructs(), books, nil
}

func (db *DB) SerieGet(id, user int) (*utils.Serie, *Books, error) {
	var qs querySerie
	if user == 0 {
		qs = querySerie{
			SelectSerie, []interface{}{id},
			func(s *Serie, b *Book, a *Author) []interface{} {
				return []interface{}{
					&b.Id, &b.title, &b.number, &b.description, &b.Isbn, &a.id, &a.name,
					&s.name,
				}
			},
		}
	} else {
		qs = querySerie{
			SelectSerieU, []interface{}{id, user},
			func(s *Serie, b *Book, a *Author) []interface{} {
				return []interface{}{
					&b.Id, &b.title, &b.number, &b.description, &b.Isbn, &a.id, &a.name,
					&s.name, &b.owned,
				}
			},
		}
	}
	return dedupSerie(db, qs)
}
