package main

import (
	"github.com/labstack/echo/middleware"
	"github.com/paul-bismuth/library/app"
)

func main() {
	rulz := app.New("./config/api.toml")

	// Middleware
	rulz.Use(middleware.Logger())
	rulz.Use(middleware.Recover())

	/* --------------------------------- API --------------------------------- */
	rulz.Api.Use(app.ContentType)
	rulz.Api.Use(middleware.CORS())

	rulz.Api.GET("/books/:isbn", rulz.Handler(app.APIBookGet), rulz.BasicAuth(false))
	rulz.Api.GET("/books/", rulz.Handler(app.APIBookList), rulz.BasicAuth(false))

	rulz.Api.POST("/books/", rulz.Handler(app.APIBookPost))

	rulz.Api.GET("/series/:id", rulz.Handler(app.APISerieGet), rulz.BasicAuth(false))
	rulz.Api.GET("/series/", rulz.Handler(app.APISerieList), rulz.BasicAuth(false))

	/* --------------------------------- WEB --------------------------------- */
	rulz.Web.Static("/static", rulz.Configuration.Paths.Static)
	rulz.Web.Static("/thumbs", rulz.Configuration.Paths.Thumbs)

	rulz.Web.GET("/", rulz.Handler(app.WEBIndex))
	rulz.Web.GET("/books/:isbn", rulz.Handler(app.WEBBookGet)).Name = "books"

	// Start application
	rulz.Logger.Fatal(rulz.Start())
}
