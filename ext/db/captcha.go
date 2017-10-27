package db

import (
	"database/sql"
)

const insertIsbn = `INSERT INTO captcha (isbn) VALUES ($1) ON CONFLICT DO NOTHING`

const removeIsbn = `DELETE FROM captcha WHERE isbn = $1`

const listIsbn = `SELECT isbn FROM captcha`

func (db *DB) CaptchaAdd(isbn string) (int64, error) {
	return db.Exec(insertIsbn, isbn)
}

func (db *DB) CaptchaRemove(isbn string) (int64, error) {
	return db.Exec(removeIsbn, isbn)
}

func (db *DB) CaptchaList() (isbns []string, err error) {
	var isbn string
	var rows *sql.Rows

	if rows, err = db.Query(listIsbn); err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&isbn); err != nil {
			return
		}
		isbns = append(isbns, isbn)
	}
	return
}
