package app

import (
	"fmt"
	"github.com/RulzUrLibrary/api/ext/auth"
	"github.com/RulzUrLibrary/api/ext/db"
	"github.com/RulzUrLibrary/api/ext/scrapper"
	"github.com/RulzUrLibrary/api/ext/smtp"
	"github.com/RulzUrLibrary/api/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"strings"
)

var ErrNotFound = echo.ErrNotFound

type dict = utils.Dict

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

func New(init Initializer) *Application {
	var app = &Application{Echo: echo.New(), Api: echo.New(), Web: echo.New()}

	app.Configuration = init.Config()

	app.Debug = app.Configuration.Debug
	app.Api.Debug = app.Configuration.Debug
	app.Web.Debug = app.Configuration.Debug

	app.Logger = init.Logger(PREFIX)
	app.Api.Logger = init.Logger(PREFIX)
	app.Web.Logger = init.Logger(PREFIX)

	app.Api.Validator = init.Validator()
	app.Web.Validator = init.Validator()

	app.Web.HTTPErrorHandler = HTTPErrorHandler
	app.Web.Renderer = init.View(app)

	app.Scrapper = init.Scrapper()
	app.Database, app.Auth = init.DB()
	app.Smtp = init.Smtp()

	if !app.Configuration.Dev {
		app.HideBanner = true
	}

	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())
	app.Use(middleware.Secure())

	app.Any("/*", func(c echo.Context) (err error) {
		if strings.HasPrefix(c.Request().Host, "api.") {
			app.Api.ServeHTTP(c.Response(), c.Request())
		} else {
			app.Web.ServeHTTP(c.Response(), c.Request())
		}
		return
	})
	/* --------------------------------- API --------------------------------- */
	app.Api.Use(ContentType)
	app.Api.Use(middleware.CORS())

	app.Api.Static("/thumbs", app.Configuration.Paths.Thumbs)

	app.Api.GET("/books/:isbn", app.Handler(APIBookGet), app.BasicAuth(false))
	app.Api.GET("/books/", app.Handler(APIBookList), app.BasicAuth(false))

	app.Api.POST("/books/", app.Handler(APIBookPost))
	app.Api.PUT("/books/", app.Handler(APIBookPut), app.BasicAuth(true))
	app.Api.DELETE("/books/", app.Handler(APIBookDelete), app.BasicAuth(true))

	app.Api.GET("/series/:id", app.Handler(APISerieGet), app.BasicAuth(false))
	app.Api.GET("/series/", app.Handler(APISerieList), app.BasicAuth(false))

	app.Api.GET("/wishlists/", app.Handler(APIWishlist), app.BasicAuth(true))

	/* --------------------------------- WEB --------------------------------- */
	app.Web.Use(CookieAuth(app.Configuration.Dev))
	app.Web.Use(I18n(app.Configuration.I18n))
	app.Web.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:X-CSRF-Token",
	}))

	app.Web.Static("/static", app.Configuration.Paths.Static)
	app.Web.Static("/thumbs", app.Configuration.Paths.Thumbs)

	app.Web.GET("/", app.Handler(WEBIndex)).Name = "index"

	app.Web.GET("/user", app.Handler(WEBUserGet), Protected).Name = "user"
	app.Web.GET("/user/activate/:id", app.Handler(WEBUserActivate)).Name = "activate"

	app.Web.POST("/user/change", app.Handler(WEBUserChange), Protected).Name = "change"

	app.Web.GET("/user/reset", app.Handler(WEBUserResetGet)).Name = "reset"
	app.Web.POST("/user/reset", app.Handler(WEBUserResetPost))

	app.Web.GET("/user/reset/:id", app.Handler(WEBUserReinit)).Name = "reinit"
	app.Web.POST("/user/reset/:id", app.Handler(WEBUserReinit))

	app.Web.POST("/user/logout", app.Handler(WEBUserLogout), Protected).Name = "logout"
	app.Web.POST("/user/lang", app.Handler(WEBUserLang)).Name = "lang"

	app.Web.GET("/books/", app.Handler(WEBBookList), Protected).Name = "books"

	app.Web.GET("/wishlists/", app.Handler(WEBWishlist), Protected).Name = "wishlists"

	app.Web.GET("/tags/", app.Handler(WEBTag), Protected).Name = "tags"
	app.Web.POST("/tags/", app.Handler(WEBTag), Protected)

	app.Web.GET("/wishlist/:id", app.Handler(WEBWishListGet)).Name = "wishlist"
	app.Web.POST("/wishlist/:id", app.Handler(WEBWishListPost), Protected)

	app.Web.GET("/books/:isbn", app.Handler(WEBBookGet)).Name = "book"
	app.Web.POST("/books/:isbn", app.Handler(WEBBookPost))

	app.Web.GET("/books/:isbn/wishlist", app.Handler(WEBWishlistAdd), Protected).Name = "share"
	app.Web.POST("/books/:isbn/wishlist", app.Handler(WEBWishlistPost), Protected)

	// TODO: find a better way to identify series
	app.Web.GET("/series/:id", app.Handler(WEBSerieGet), Protected).Name = "serie"
	app.Web.POST("/series/:id", app.Handler(WEBSeriePost), Protected)

	app.Web.GET("/auth", app.Handler(WEBAuthGet)).Name = "auth"
	app.Web.POST("/auth", app.Handler(WEBAuthPost))

	app.Web.GET("/auth/new", app.Handler(WEBUserNewGet)).Name = "new"
	app.Web.POST("/auth/new", app.Handler(WEBUserNewPost))

	return app
}

func (app *Application) Start() error {
	host := app.Configuration.Host
	port := app.Configuration.Port

	if !app.Configuration.Dev {
		go func() {
			e := echo.New()
			e.Pre(middleware.HTTPSRedirectWithConfig(middleware.RedirectConfig{
				Code: http.StatusTemporaryRedirect,
			}))
			e.Start(":80")
		}()
		app.Echo.AutoTLSManager.Cache = autocert.DirCache(app.Configuration.Paths.TLSCache)
		return app.Echo.StartAutoTLS(":443")
	} else {
		return app.Echo.Start(fmt.Sprintf("%s:%d", host, port))
	}
}
