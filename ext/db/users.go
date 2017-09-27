package db

import (
	"database/sql"
	"github.com/ixday/echo-hello/utils"
)

const authUser = `
SELECT COALESCE(pwhash = crypt($2, pwhash), FALSE), id
FROM users
WHERE name = $1`

const authGoogle = `
WITH s AS (
	SELECT id FROM users WHERE name = $1
), i AS (
	INSERT INTO users ("name") SELECT $1 WHERE NOT EXISTS (SELECT 1 FROM s) RETURNING id
)
SELECT id FROM i UNION ALL SELECT id FROM s`

const newUser = `
WITH s AS (
  SELECT id, false FROM users WHERE name = $1
), i AS (
  INSERT INTO users ("name", "pwhash") SELECT $1, crypt($2, gen_salt('bf')) WHERE NOT EXISTS (SELECT 1 FROM s) RETURNING id, true
)
SELECT id, bool FROM i UNION ALL SELECT id, bool FROM s`

func (db *DB) Auth(name, password string) (u utils.User, err error) {
	var ok bool

	if err = db.QueryRow(authUser, name, password).Scan(&ok, &u.Id); err == nil {
		if ok {
			u.Name = name
		} else {
			err = utils.ErrUserAuth
		}
	} else {
		if err == sql.ErrNoRows {
			err = utils.ErrUserAuth
		}

	}
	return
}

func (db *DB) AuthGoogle(name string) (user utils.User, err error) {
	user.Name = name
	err = db.QueryRow(authGoogle, name).Scan(&user.Id)
	return
}

func (db *DB) NewUser(name, password string) (user utils.User, err error) {
	var ok bool
	user.Name = name
	err = db.QueryRow(newUser, name, password).Scan(&user.Id, &ok)
	if !ok {
		return user, utils.ErrUserExists
	}
	return
}
