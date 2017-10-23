package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/rulzurlibrary/api/app"
)

func main() {
	rulz := app.New()

	// Middleware
	rulz.Use(middleware.Logger())
	rulz.Use(middleware.Recover())
	rulz.Use(middleware.Secure())

	/* --------------------------------- API --------------------------------- */
	rulz.Api.Use(app.ContentType)
	rulz.Api.Use(middleware.CORS())

	rulz.Api.GET("/books/:isbn", rulz.Handler(app.APIBookGet), rulz.BasicAuth(false))
	rulz.Api.GET("/books/", rulz.Handler(app.APIBookList), rulz.BasicAuth(false))

	rulz.Api.POST("/books/", rulz.Handler(app.APIBookPost))
	rulz.Api.PUT("/books/", rulz.Handler(app.APIBookPut), rulz.BasicAuth(true))
	rulz.Api.DELETE("/books/", rulz.Handler(app.APIBookDelete), rulz.BasicAuth(true))

	rulz.Api.GET("/series/:id", rulz.Handler(app.APISerieGet), rulz.BasicAuth(false))
	rulz.Api.GET("/series/", rulz.Handler(app.APISerieList), rulz.BasicAuth(false))

	/* --------------------------------- WEB --------------------------------- */
	rulz.Web.Use(app.CookieAuth(rulz.Configuration.Dev))
	rulz.Web.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:X-CSRF-Token",
	}))

	rulz.Web.Static("/static", rulz.Configuration.Paths.Static)
	rulz.Web.Static("/thumbs", rulz.Configuration.Paths.Thumbs)

	rulz.Web.GET("/", rulz.Handler(app.WEBIndex)).Name = "index"

	rulz.Web.GET("/user", rulz.Handler(app.WEBUserGet), app.Protected).Name = "user"
	rulz.Web.GET("/user/reset", rulz.Handler(app.WEBUserReset), app.Protected).Name = "reset"

	rulz.Web.POST("/user/password", rulz.Handler(app.WEBUserPassword), app.Protected).Name = "password"
	rulz.Web.POST("/user/logout", rulz.Handler(app.WEBUserLogout), app.Protected).Name = "logout"
	rulz.Web.POST("/user/lang", rulz.Handler(app.WEBUserLang)).Name = "lang"

	rulz.Web.GET("/books/", rulz.Handler(app.WEBBookList), app.Protected).Name = "books"

	rulz.Web.GET("/books/:isbn", rulz.Handler(app.WEBBookGet)).Name = "book"
	rulz.Web.POST("/books/:isbn", rulz.Handler(app.WEBBookPost))

	// TODO: find a better way to identify series
	rulz.Web.GET("/series/:id", rulz.Handler(app.WEBSerieGet), app.Protected).Name = "serie"
	rulz.Web.POST("/series/:id", rulz.Handler(app.WEBSeriePost), app.Protected)

	rulz.Web.GET("/auth", rulz.Handler(app.WEBAuthGet)).Name = "auth"
	rulz.Web.POST("/auth", rulz.Handler(app.WEBAuthPost))

	rulz.Web.GET("/auth/new", rulz.Handler(app.WEBUserNewGet)).Name = "new"
	rulz.Web.POST("/auth/new", rulz.Handler(app.WEBUserNewPost))

	// Start application
	rulz.Logger.Fatal(rulz.Start())
}
