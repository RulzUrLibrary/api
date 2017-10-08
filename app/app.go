package app

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/paul-bismuth/library/ext/auth"
	"github.com/paul-bismuth/library/ext/db"
	"github.com/paul-bismuth/library/ext/scrapper"
	"github.com/paul-bismuth/library/ext/view"
	"net/http"
	"strings"
)

type Context struct {
	echo.Context
	DB       *db.DB
	Auth     *auth.Auth
	Logger   echo.Logger
	Scrapper *scrapper.Scrapper
}

type Application struct {
	*echo.Echo
	Api           *echo.Echo
	Web           *echo.Echo
	Database      *db.DB
	Scrapper      *scrapper.Scrapper
	Auth          *auth.Auth
	Configuration Configuration
}

func (app *Application) Handler(h func(*Context) error) echo.HandlerFunc {
	return func(original echo.Context) error {
		return h(&Context{original, app.Database, app.Auth, app.Logger, app.Scrapper})
	}
}

func New(configPath string) *Application {
	var err error
	var app = &Application{Echo: echo.New(), Api: echo.New(), Web: echo.New()}

	app.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		if strings.HasPrefix(req.Host, "api.") {
			app.Api.ServeHTTP(res, req)
		} else {
			app.Web.ServeHTTP(res, req)
		}
		return
	})

	app.Configuration, err = ParseConfig(configPath)
	if err != nil {
		app.Logger.Fatal(err)
	}

	app.Debug = app.Configuration.Debug
	app.Web.Renderer = view.New(view.Configuration{
		app.Configuration.Paths.Templates,
		app.Configuration.Dev,
		app.Web,
	})
	app.Scrapper = scrapper.New(app.Logger, app.Configuration.Paths.Thumbs)
	app.Database = db.New(app.Logger, app.Configuration.Database)
	app.Auth = auth.New(app.Logger, app.Database)

	if !app.Configuration.Dev {
		app.HideBanner = true
	}
	return app
}

func (app *Application) Start() error {
	host := app.Configuration.Host
	port := app.Configuration.Port

	if !app.Configuration.Dev {
		key := app.Configuration.Paths.Key
		cert := app.Configuration.Paths.Cert
		go func() {
			e := echo.New()
			e.Pre(middleware.HTTPSRedirectWithConfig(middleware.RedirectConfig{
				Code: http.StatusTemporaryRedirect,
			}))
			e.Start(":80")
		}()
		return app.Echo.StartTLS(":443", key, cert)
	} else {

		return app.Echo.Start(fmt.Sprintf("%s:%d", host, port))
	}
}
