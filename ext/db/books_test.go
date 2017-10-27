package db

import (
	"bytes"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"github.com/rulzurlibrary/api/ext/db"
	"os"
	"path/filepath"
	"testing"
)

const query = `
SELECT array_agg(DISTINCT ROW(a.id, a.name))
FROM books b
INNER JOIN series s ON (b.fk_serie = s.id)
LEFT OUTER JOIN collections c ON (b.id = fk_book AND fk_user = 1)
LEFT OUTER JOIN book_authors ba ON (b.id = ba.fk_book)
LEFT OUTER JOIN authors a ON (ba.fk_author = a.id)
LEFT OUTER JOIN wishlists_books wb ON (b.id = wb.fk_book)
LEFT OUTER JOIN wishlists w ON (w.id = wb.fk_wishlist)
WHERE b.isbn = '9782800100173'
GROUP BY b.id, s.name, c.fk_book`

const CONFIG_ENV = "RULZURLIBRARY_CONFIG"

func NewDB() (*db.DB, error) {
	configuration := struct {
		Database db.Configuration
	}{}

	filename, err := filepath.Abs(os.Getenv(CONFIG_ENV))
	if err != nil {
		return nil, err
	}

	// put the file's contents as toml to the default configuration(c)
	_, err = toml.DecodeFile(filename, &configuration)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	logger := log.New("testing")
	logger.SetLevel(log.OFF)
	return db.New(logger, configuration.Database), nil
}

func TestArrayScan(t *testing.T) {
	db, err := NewDB()
	if err != nil {
		t.Errorf("failing init db %s", err)
	}
	authors := []db.Author{}
	if err := db.QueryRow(query).Scan(pq.Array(&authors)); err != nil {
		t.Errorf("failing scanning %s", err)
	}
	fmt.Printf("%#v\n", authors)
}
