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
SELECT b.id, title, num, isbn, s.id, s.name, a.id, a.name,
	(SELECT b.description WHERE b.num IS NULL OR b.num = 1),
	EXISTS(SELECT true FROM collections WHERE fk_book = b.id AND fk_user = $3)
FROM (
    SELECT s.id FROM books b, collections, series s
    WHERE fk_user = $3 AND fk_book = b.id AND fk_serie = s.id
    GROUP BY s.id ORDER BY s.name DESC NULLS LAST LIMIT $1 OFFSET $2
) r, books b
INNER JOIN series s ON (fk_serie = s.id)
LEFT OUTER JOIN book_authors ba ON (b.id = fk_book)
LEFT JOIN authors a ON (fk_author = a.id)
WHERE s.id = r.id ORDER BY s.name, num`

const SelectSeries = `
SELECT b.id, title, num, isbn, s.id, s.name, a.id, a.name,
	(SELECT b.description WHERE b.num IS NULL OR b.num = 1)
FROM (
  SELECT id FROM series
  GROUP BY id ORDER BY id DESC NULLS LAST LIMIT $1 OFFSET $2
) r, books b
INNER JOIN series s ON (fk_serie = s.id)
LEFT OUTER JOIN book_authors ba ON (b.id = fk_book)
LEFT JOIN authors a ON (fk_author = a.id)
WHERE s.id = r.id ORDER BY s.id, num`

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

type querySerieList struct {
	querySerie
	queryList     string
	queryListArgs []interface{}
}

func dedupSeries(db *DB, qs querySerie) (*Series, error) {
	series := NewSeries()

	rows, err := db.Query(qs.query, qs.args...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var serie Serie
		var author Author
		var volume Book

		if err = rows.Scan(qs.getArgs(&serie, &volume, &author)...); err != nil {
			return nil, err
		}

		if s := series.Get(serie.id); s != nil {
			s.authors = append(s.authors, author)
			if v := s.volumes.Get(volume.Id); v != nil {
				v.authors = append(v.authors, author)
			} else {
				volume.authors = Authors{author}
				s.volumes.Add(&volume)
			}
		} else {
			serie.authors = Authors{author}
			volume.authors = Authors{author}

			serie.volumes = NewBooks()
			serie.volumes.Add(&volume)
			series.Add(&serie)
		}
	}
	return series, nil
}

func (db *DB) SerieGet(id int) (*Serie, error) {
	series, err := dedupSeries(db, querySerie{
		SelectSerie, []interface{}{id},
		func(s *Serie, b *Book, a *Author) []interface{} {
			return []interface{}{
				&b.Id, &b.title, &b.number, &b.description, &b.Isbn, &a.id, &a.name,
				&s.name,
			}
		},
	})
	if err != nil {
		return nil, err
	}
	return series.First(), nil
}

func (db *DB) SerieGetU(id, user int) (*Serie, error) {
	series, err := dedupSeries(db, querySerie{
		SelectSerieU, []interface{}{id, user},
		func(s *Serie, b *Book, a *Author) []interface{} {
			return []interface{}{
				&b.Id, &b.title, &b.number, &b.description, &b.Isbn, &a.id, &a.name,
				&s.name, &b.owned,
			}
		},
	})
	if err != nil {
		return nil, err
	}
	return series.First(), nil
}

func (db *DB) serieList(qsl querySerieList) (*Series, int, error) {
	count, err := db.Count(qsl.queryList, qsl.queryListArgs...)
	if err != nil {
		return nil, 0, err
	}
	series, err := dedupSeries(db, qsl.querySerie)
	if err != nil {
		return nil, 0, err
	}
	return series, count, nil

}

func (db *DB) SerieList(limit, offset int) (*Series, int, error) {
	return db.serieList(querySerieList{
		querySerie{
			SelectSeries, []interface{}{limit, offset},
			func(s *Serie, b *Book, a *Author) []interface{} {
				return []interface{}{
					&b.Id, &b.title, &b.number, &b.Isbn, &s.id, &s.name, &a.id, &a.name,
					&s.description,
				}
			},
		}, CountSeries, []interface{}{},
	})
}

func (db *DB) SerieListU(limit, offset, user int) (*Series, int, error) {
	return db.serieList(querySerieList{
		querySerie{
			SelectSeriesU, []interface{}{limit, offset, user},
			func(s *Serie, b *Book, a *Author) []interface{} {
				return []interface{}{
					&b.Id, &b.title, &b.number, &b.Isbn, &s.id, &s.name, &a.id, &a.name,
					&s.description, &b.owned,
				}
			},
		}, CountSeriesU, []interface{}{user},
	})
}
