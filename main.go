package main

import (
	"github.com/ixday/echo-hello/app"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
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
	rulz.Api.POST("/books/", rulz.Handler(app.APIBookPost))

	/* --------------------------------- WEB --------------------------------- */
	rulz.Web.Static("/static", rulz.Configuration.Paths.Static)

	rulz.Web.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"name": "Dolly!",
		})
	})

	// Start application
	rulz.Logger.Fatal(rulz.Start())
}
