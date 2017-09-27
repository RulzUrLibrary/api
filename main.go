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
	rulz.Api.Use(middleware.CORS())
	rulz.Api.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	rulz.Api.GET("/books/:isbn",
		rulz.Handler(func(c *app.Context) error {
			book, err := app.BookGet(c)
			if err == nil {
				c.JSON(http.StatusOK, book)
			}
			return err
		}),
		rulz.BasicAuth(false),
	)

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
