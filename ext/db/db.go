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
