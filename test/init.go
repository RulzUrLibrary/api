package test

import (
	"github.com/labstack/gommon/log"
	"github.com/rulzurlibrary/api/app"
	"github.com/rulzurlibrary/api/ext/db"
	fakeDB "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io"
	"net/http"
	"net/http/httptest"
)

const (
	PREFIX = "rulz_test"
)

var (
	rulz   *app.Application
	logger *log.Logger
	mock   fakeDB.Sqlmock
)

func initDB() *db.DB {
	var err error

	fake := &db.DB{Logger: logger}
	if fake.DB, mock, err = fakeDB.New(); err != nil {
		logger.Fatal(err)
	}
	return fake
}

func NewRequestAPI(method, target string, body io.Reader) (req *http.Request) {
	req = httptest.NewRequest(method, target, body)
	req.Host = "api."
	return
}

func init() {
	logger = log.New(PREFIX)
	logger.SetLevel(log.DEBUG)

	rulz = app.New(initDB(), app.Configuration{Debug: true})
}
