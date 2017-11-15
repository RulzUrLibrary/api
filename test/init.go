package test

import (
	"github.com/labstack/gommon/log"
	"github.com/rulzurlibrary/api/app"
	"github.com/rulzurlibrary/api/ext/auth"
	"github.com/rulzurlibrary/api/ext/db"
	fakeDB "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io"
	"net/http"
	"net/http/httptest"
)

var (
	rulz *app.Application
	mock fakeDB.Sqlmock
)

type TestInitializer struct {
	*app.DefaultInitializer
}

func NewTestInitializer() *TestInitializer {
	return &TestInitializer{
		&app.DefaultInitializer{app.Configuration{Debug: true}, log.DEBUG},
	}
}

func (ti *TestInitializer) Logger(prefix string) *log.Logger {
	return ti.DefaultInitializer.Logger("test_" + prefix)
}

func (ti *TestInitializer) DB() (*db.DB, *auth.Auth) {
	fake, mocker, err := fakeDB.New()

	if err != nil {
		ti.Logger(app.PREFIX).Fatal(err)
	}
	mock = mocker

	database := &db.DB{fake, ti.Logger(app.PREFIX_DB)}
	auth := auth.New(ti.Logger(app.PREFIX_AUTH), database)
	return database, auth
}

func NewRequestAPI(method, target string, body io.Reader) (req *http.Request) {
	req = httptest.NewRequest(method, target, body)
	req.Host = "api."
	return
}

func init() {
	rulz = app.New(NewTestInitializer())
}
