package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
)

type Configuration struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     int
}

func (c Configuration) String() string {
	return fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%d sslmode=disable",
		c.Name, c.User, c.Password, c.Host, c.Port,
	)
}

func (c Configuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name     string `json:"name"`
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	}{c.Name, c.User, "******", c.Host, c.Port})
}

type DB struct {
	*sql.DB
}

func New(c Configuration) (*DB, error) {
	db, err := sql.Open("postgres", c.String())
	return &DB{db}, err
}

func (db *DB) Exists(from, where string, args ...interface{}) (ok bool, err error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s=$1)", from, where)
	err = db.QueryRow(query, args...).Scan(&ok)
	return
}

func (db *DB) Count(query string, args ...interface{}) (count int, err error) {
	err = db.QueryRow(query, args...).Scan(&count)
	return
}

func toInterfaceS(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func toInterfaceI(i int) interface{} {
	if i == 0 {
		return nil
	}
	return i
}

type Tx struct {
	*sql.Tx
}

func (tx *Tx) Insert(query string, args ...interface{}) (id int, err error) {
	//if glog.V(2) {
	//	glog.Infoln(append([]interface{}{query}, args...)...)
	//}
	err = tx.QueryRow(query, args...).Scan(&id)
	return
}

func (db *DB) Transaction(clojure func(*Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = clojure(&Tx{tx})
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err == nil {
		return nil
	}
	// log error
	if err := tx.Rollback(); err != nil {
		return err
	}
	return err
}
