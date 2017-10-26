package app

import (
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/rulzurlibrary/api/ext/auth"
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/ext/scrapper"
	"github.com/rulzurlibrary/api/ext/smtp"
	"github.com/rulzurlibrary/api/ext/validator"
	"github.com/rulzurlibrary/api/ext/view"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
	"strings"
)

type dict = utils.Dict

// aliases
var ErrNotFound = echo.ErrNotFound

func dynamic(c *Context, from, where, param string) error {
	ok, err := c.App.Database.Exists(from, where, param)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNotFound
	}
	return nil
}

type Context struct {
	echo.Context
	App *Application
}

func (c *Context) Flashes(flashes ...utils.Flash) error {
	session, _ := c.Get("session").(*sessions.Session)
	for _, flash := range flashes {
		session.AddFlash(flash)
	}
	return session.Save(c.Request(), c.Response())
}

func (c *Context) Reverse(name string, params ...interface{}) string {
	return c.Echo().Reverse(name, params...)
}

func (c *Context) ReverseAbs(name string, params ...interface{}) string {
	host := c.App.Configuration.Host
	port := c.App.Configuration.Port

	if c.App.Configuration.Dev {
		return fmt.Sprintf("http://%s:%d%s", host, port, c.Reverse(name, params...))
	} else {
		return fmt.Sprintf("https://%s%s", host, c.Reverse(name, params...))
	}
}

func (c *Context) RedirectWithFlash(msg string) error {
	flash := utils.Flash{utils.FlashSuccess, msg}
	if err := c.Flashes(flash); err != nil {
		return err
	}
	redirect := "index"
	if _, ok := c.Get("user").(*utils.User); ok {
		redirect = "books"
	}
	return c.Redirect(http.StatusSeeOther, c.Echo().Reverse(redirect))
}

func (c *Context) SaveUser(user *utils.User, flashes ...utils.Flash) error {
	session, _ := c.Get("session").(*sessions.Session)
	if user == nil {
		session.Values["user"] = nil
	} else {
		session.Values["user"] = user
	}
	c.Set("user", user)
	for _, flash := range flashes {
		session.AddFlash(flash)
	}
	return session.Save(c.Request(), c.Response())
}

type Application struct {
	*echo.Echo
	Api           *echo.Echo
	Web           *echo.Echo
	Database      *db.DB
	Scrapper      *scrapper.Scrapper
	Auth          *auth.Auth
	Smtp          *smtp.Smtp
	Configuration Configuration
}

func (app *Application) Handler(h func(*Context) error) echo.HandlerFunc {
	return func(original echo.Context) error {
		return h(&Context{original, app})
	}
}

func HTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	data := dict{"error": err, "msg": "Internal server error"}

	if httpErr, ok := err.(*echo.HTTPError); ok {
		code = httpErr.Code
		data["err"] = httpErr.Inner
		data["msg"] = httpErr.Message
	}
	data["code"] = code
	err = c.Render(code, "error.html", data)
	if err != nil {
		c.Logger().Error(err)
	}
}

func New() *Application {
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

	app.Configuration, err = ParseConfig()
	if err != nil {
		app.Logger.Fatal(err)
	}

	app.Api.Debug = app.Configuration.Debug
	app.Api.Validator = validator.New()

	app.Web.Debug = app.Configuration.Debug
	app.Web.Validator = validator.New()
	app.Web.HTTPErrorHandler = HTTPErrorHandler
	app.Web.Renderer = view.New(app.Web, app.Configuration.View)

	app.Debug = app.Configuration.Debug
	app.Scrapper = scrapper.New(app.Logger, app.Configuration.Paths.Thumbs)
	app.Database = db.New(app.Logger, app.Configuration.Database)
	app.Auth = auth.New(app.Logger, app.Database)
	app.Smtp = smtp.New(app.Logger, app.Configuration.Smtp)

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
		return app.Echo.StartTLS(":443", cert, key)
	} else {

		return app.Echo.Start(fmt.Sprintf("%s:%d", host, port))
	}
}
