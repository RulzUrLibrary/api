package db

import (
	"database/sql"
	"github.com/rulzurlibrary/api/utils"
)

const authUser = `
SELECT COALESCE(pwhash = crypt($2, pwhash), FALSE), id
FROM users
WHERE name = $1`

const changePassword = `
UPDATE users SET pwhash = crypt($1, gen_salt('bf'))
WHERE COALESCE(pwhash = crypt($2, pwhash), FALSE) AND id = $3
`

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

func (db *DB) Auth(name, password string) (*utils.User, error) {
	var ok bool
	var user = &utils.User{Name: name}

	err := db.QueryRow(authUser, name, password).Scan(&ok, &user.Id)

	if err != nil && err == sql.ErrNoRows || !ok {
		return nil, utils.ErrUserAuth
	}
	return user, err
}

func (db *DB) AuthGoogle(name string) (*utils.User, error) {
	user := &utils.User{Name: name}
	return user, db.QueryRow(authGoogle, name).Scan(&user.Id)
}

func (db *DB) ChangePassword(new, old string, user int) (int, error) {
	return db.Exec(changePassword, new, old, user)
}

func (db *DB) NewUser(name, password string) (*utils.User, error) {
	var ok bool
	var user = &utils.User{Name: name}

	err := db.QueryRow(newUser, name, password).Scan(&user.Id, &ok)

	if err == nil && !ok {
		return nil, utils.ErrUserExists
	}
	return user, err
}
