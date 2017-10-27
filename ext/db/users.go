package db

import (
	"database/sql"
	"github.com/rulzurlibrary/api/utils"
)

const authUser = `
SELECT COALESCE(pwhash = crypt($2, pwhash), FALSE), id
FROM users
WHERE email = $1`

const passwordChange = `
UPDATE users SET pwhash = crypt($1, gen_salt('bf'))
WHERE COALESCE(pwhash = crypt($2, pwhash), FALSE) AND id = $3
`

const passwordReset = `
UPDATE users SET pwhash = crypt($1, gen_salt('bf')), reset = null WHERE reset = $2`

const authGoogle = `
WITH s AS (
	SELECT id FROM users WHERE email = $1
), i AS (
	INSERT INTO users ("email") SELECT $1 WHERE NOT EXISTS (SELECT 1 FROM s) RETURNING id
)
SELECT id FROM i UNION ALL SELECT id FROM s`

const newUser = `
WITH s AS (
  SELECT id, activate, false FROM users WHERE email = $1
), i AS (
  INSERT INTO users ("email", "pwhash", "activate")
  SELECT $1, crypt($2, gen_salt('bf')), gen_random_uuid()
  WHERE NOT EXISTS (SELECT 1 FROM s)
  RETURNING id, activate, true
)
SELECT id, activate, bool FROM i UNION ALL SELECT id, activate, bool FROM s`

const deleteUser = `
DELETE FROM users WHERE id = $1`

const deleteActivate = `
UPDATE users SET activate = null WHERE activate = $1`

const createReset = `
UPDATE users SET reset = gen_random_uuid() WHERE email = $1 RETURNING reset`

func (db *DB) Auth(email, password string) (*utils.User, error) {
	var ok bool
	var user = &utils.User{Email: email}

	err := db.QueryRow(authUser, email, password).Scan(&ok, &user.Id)

	if err != nil && err == sql.ErrNoRows || !ok {
		return nil, utils.ErrUserAuth
	}
	return user, err
}

func (db *DB) AuthGoogle(email string) (*utils.User, error) {
	user := &utils.User{Email: email}
	return user, db.QueryRow(authGoogle, email).Scan(&user.Id)
}

func (db *DB) PasswordChange(new, old string, user int) (int64, error) {
	return db.Exec(passwordChange, new, old, user)
}

func (db *DB) PasswordReset(new, reset string) error {
	_, err := db.DB.Exec(passwordReset, new, reset)
	return err
}

func (db *DB) NewUser(email, password string) (*utils.User, string, error) {
	var ok bool
	var activate sql.NullString
	var user = &utils.User{Email: email}
	err := db.QueryRow(newUser, email, password).Scan(&user.Id, &activate, &ok)

	if err == nil && !ok {
		return nil, "", utils.ErrUserExists
	}
	return user, activate.String, err
}

func (db *DB) DeleteUser(user *utils.User) (int64, error) {
	return db.Exec(deleteUser, user.Id)
}

func (db *DB) DeleteActivate(activate string) error {
	count, err := db.Exec(deleteActivate, activate)
	if err != nil {
		return err
	}
	if count == 0 {
		return utils.ErrAlreadyActivate
	}
	return nil
}

func (db *DB) MustDeleteUser(user *utils.User) {
	count, err := db.DeleteUser(user)
	if count != 1 {
		panic("no user removed")
	}
	if err != nil {
		panic(err)
	}
}

func (db *DB) CreateReset(email string) (reset string, err error) {
	// we don't consider user non existing as an error,
	// to avoid email registered guessing
	err = db.QueryRow(createReset, email).Scan(&reset)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}
