package db

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"strings"
)

type list = []interface{}

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
	Logger echo.Logger
}

func New(l echo.Logger, c Configuration) *DB {
	db, err := sql.Open("postgres", c.String())
	if err != nil {
		l.Fatal(err)
	}
	return &DB{db, l}
}

func (db *DB) Exists(from, where string, args ...interface{}) (ok bool, err error) {
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE %s=$1)", from, where)
	err = db.QueryRow(query, args...).Scan(&ok)
	return
}

func (db *DB) Count(query string, args ...interface{}) (count int64, err error) {
	err = db.QueryRow(query, args...).Scan(&count)
	return
}

func (db *DB) Exec(query string, args ...interface{}) (int64, error) {
	res, err := db.DB.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
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

func (db *DB) Transaction(clojure func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	err = clojure(tx)
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

func (db *DB) query(query string, args list, scan func() list) error {
	rows, err := db.Query(query, args...)
	if err != nil {
		return err
	}

	for rows.Next() {
		if err := rows.Scan(scan()...); err != nil {
			return err
		}
	}
	return err
}

func (db *DB) queryList(query string, args list, scan func() list, queryC string, argsC list) (int64, error) {
	count, err := db.Count(queryC, argsC...)
	if err != nil {
		return 0, err
	}

	if err := db.query(query, args, scan); err != nil {
		return 0, err
	}
	return count, nil
}

func emptyArray(size int, src interface{}) bool {
	return bytes.Equal(src.([]byte), []byte(fmt.Sprintf(`{"(%s)"}`, strings.Repeat(",", size-1))))
}

// Come from https://github.com/lib/pq/blob/b609790bd85edf8e9ab7e0f8912750a786177bcf/array.go#L642
func parseRow(src, del []byte) (elems [][]byte, err error) {
	var depth, i int
	var dims []int

	if len(src) < 1 || src[0] != '(' {
		return nil, fmt.Errorf("pq: unable to parse array; expected %q at offset %d", '(', 0)
	}

Open:
	for i < len(src) {
		switch src[i] {
		case '(':
			depth++
			i++
		case ')':
			elems = make([][]byte, 0)
			goto Close
		default:
			break Open
		}
	}
	dims = make([]int, i)

Element:
	for i < len(src) {
		switch src[i] {
		case '(':
			if depth == len(dims) {
				break Element
			}
			depth++
			dims[depth-1] = 0
			i++
		case '"':
			var elem = []byte{}
			var escape bool
			for i++; i < len(src); i++ {
				if escape {
					elem = append(elem, src[i])
					escape = false
				} else {
					switch src[i] {
					default:
						elem = append(elem, src[i])
					case '\\':
						escape = true
					case '"':
						elems = append(elems, elem)
						i++
						break Element
					}
				}
			}
		default:
			for start := i; i < len(src); i++ {
				if bytes.HasPrefix(src[i:], del) || src[i] == ')' {
					elem := src[start:i]
					if len(elem) == 0 {
						return make([][]byte, 0), nil
					}
					if bytes.Equal(elem, []byte("NULL")) {
						elem = nil
					}
					elems = append(elems, elem)
					break Element
				}
			}
		}
	}

	for i < len(src) {
		if bytes.HasPrefix(src[i:], del) && depth > 0 {
			dims[depth-1]++
			i += len(del)
			goto Element
		} else if src[i] == ')' && depth > 0 {
			dims[depth-1]++
			depth--
			i++
		} else {
			return nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}

Close:
	for i < len(src) {
		if src[i] == ')' && depth > 0 {
			depth--
			i++
		} else {
			return nil, fmt.Errorf("pq: unable to parse array; unexpected %q at offset %d", src[i], i)
		}
	}
	if depth > 0 {
		err = fmt.Errorf("pq: unable to parse array; expected %q at offset %d", ')', i)
	}
	if err == nil {
		for _, d := range dims {
			if (len(elems) % d) != 0 {
				err = fmt.Errorf("pq: multidimensional arrays must have elements with matching dimensions")
			}
		}
	}
	return
}
