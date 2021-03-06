package db

const CountSeries = `
SELECT COUNT(DISTINCT(id))
FROM series`

const CountSeriesU = `
SELECT COUNT(DISTINCT(id))
FROM (
	SELECT s.id FROM books b, collections, series s
	WHERE fk_user = $1 AND fk_book = b.id AND b.fk_serie = s.id
) _`

const SelectSeriesU = `
SELECT b.id, title, num, isbn, s.id, s.name,
  array_agg(DISTINCT ROW(a.id, a.name)), c.fk_book IS NOT NULL
FROM (
    SELECT s.id FROM books b, collections, series s
    WHERE fk_user = $3 AND fk_book = b.id AND fk_serie = s.id
    GROUP BY s.id ORDER BY s.name NULLS LAST LIMIT $1 OFFSET $2
) r, books b
INNER JOIN series s ON (fk_serie = s.id)
LEFT OUTER JOIN collections c ON (c.fk_book = b.id AND fk_user = $3)
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT JOIN authors a ON (fk_author = a.id)
WHERE s.id = r.id
GROUP BY b.id, b.isbn, b.title, b.description, b.price, b.num, s.id, s.name, c.fk_book
ORDER BY s.name, num`

const SelectSeries = `
SELECT b.id, title, num, isbn, s.id, s.name, array_agg(DISTINCT ROW(a.id, a.name))
FROM (
  SELECT id FROM series
  GROUP BY id ORDER BY id DESC NULLS LAST LIMIT $1 OFFSET $2
) r, books b
INNER JOIN series s ON (fk_serie = s.id)
LEFT OUTER JOIN book_authors ba ON (b.id = fk_book)
LEFT JOIN authors a ON (fk_author = a.id)
WHERE s.id = r.id GROUP BY b.id, s.id ORDER BY s.id, num`

const SelectSerie = `
SELECT b.id, title, num, description, isbn, array_agg(DISTINCT ROW(a.id, a.name)),
	s.id, s.name
FROM books b
INNER JOIN series s ON (fk_serie = s.id)
LEFT OUTER JOIN book_authors ba ON (b.id = fk_book)
LEFT JOIN authors a ON (fk_author = a.id)
WHERE s.id = $1 GROUP BY b.id, s.id ORDER BY num`

const SelectSerieU = `
SELECT b.id, title, num, description, isbn, array_agg(DISTINCT ROW(a.id, a.name)),
	s.id, s.name, c.fk_book IS NOT NULL
FROM books b
INNER JOIN series s ON (fk_serie = s.id)
LEFT OUTER JOIN collections c ON (c.fk_book = b.id AND fk_user = $2)
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT JOIN authors a ON (fk_author = a.id)
WHERE s.id = $1 GROUP BY b.id, s.id, c.fk_book ORDER BY num`

func (db *DB) SerieGet(id int) (books Books, err error) {
	err = db.query(SelectSerie, list{id}, books.InsertBook(func(b *Book) list {
		return list{&b.id, &b.title, &b.number, &b.description, &b.isbn, &b.authors, &b.serie_id, &b.serie}
	}))
	return
}

func (db *DB) SerieGetU(id, user int) (books Books, err error) {
	err = db.query(SelectSerieU, list{id, user}, books.InsertBook(func(b *Book) list {
		return list{&b.id, &b.title, &b.number, &b.description, &b.isbn, &b.authors, &b.serie_id, &b.serie, &b.owned}
	}))
	return
}

func (db *DB) SerieList(limit, offset int) (books Books, count int64, err error) {
	count, err = db.queryList(SelectSeries, list{limit, offset}, books.InsertBook(func(b *Book) list {
		return list{&b.id, &b.title, &b.number, &b.isbn, &b.serie_id, &b.serie, &b.authors}
	}), CountSeries, list{})
	return
}

func (db *DB) SerieListU(limit, offset, user int) (books Books, count int64, err error) {
	count, err = db.queryList(SelectSeriesU, list{limit, offset, user}, books.InsertBook(func(b *Book) list {
		return list{&b.id, &b.title, &b.number, &b.isbn, &b.serie_id, &b.serie, &b.authors, &b.owned}
	}), CountSeriesU, []interface{}{user})
	return
}
