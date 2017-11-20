package test

import (
	"github.com/labstack/gommon/log"
	"github.com/RulzUrLibrary/api/app"
	"github.com/RulzUrLibrary/api/ext/auth"
	"github.com/RulzUrLibrary/api/ext/db"
	"github.com/RulzUrLibrary/api/utils"
	fakeDB "gopkg.in/DATA-DOG/go-sqlmock.v1"
	"io"
	"net/http"
	"net/http/httptest"
)

var (
	rulz *app.Application
	mock fakeDB.Sqlmock
)

type (
	MockCache       struct{}
	TestInitializer struct {
		*app.DefaultInitializer
	}
)

func (mc MockCache) Get(_ string) (*utils.User, bool) {
	return &utils.User{}, false // nothing can get out from the cache
}

func (mc MockCache) Set(_ string, _ *utils.User) {}

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
	auth := auth.New(ti.Logger(app.PREFIX_AUTH), database, MockCache{})
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
