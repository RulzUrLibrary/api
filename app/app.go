package app

import (
	"fmt"
	"github.com/ixday/echo-hello/ext/db"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"html/template"
	"net/http"
	"strings"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Context struct {
	echo.Context
	*db.DB
}

type Application struct {
	*echo.Echo
	Api           *echo.Echo
	Site          *echo.Echo
	Database      *db.DB
	Configuration Configuration
}

func (app *Application) Handler(h func(*Context) error) echo.HandlerFunc {
	return func(original echo.Context) error {
		return h(&Context{original, app.Database})
	}
}

func New(configPath string) *Application {
	var err error
	var app = &Application{Echo: echo.New(), Api: echo.New(), Site: echo.New()}

	app.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		if strings.HasPrefix(req.Host, "api.") {
			app.Api.ServeHTTP(res, req)
		} else {
			app.Site.ServeHTTP(res, req)
		}
		return
	})

	app.Configuration, err = ParseConfig(configPath)
	if err != nil {
		app.Logger.Fatal(err)
	}

	app.Debug = app.Configuration.Debug

	if !app.Configuration.Dev {
		app.HideBanner = true
	}

	app.Database, err = db.New(app.Configuration.Database)
	if err != nil {
		app.Logger.Fatal(err)
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
