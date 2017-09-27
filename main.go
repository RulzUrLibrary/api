package main

import (
	"github.com/ixday/echo-hello/app"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
)

func main() {
	app := app.New("./config/api.toml")

	// Middleware
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	app.Api.Use(middleware.CORS())
	app.Api.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	app.Web.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"name": "Dolly!",
		})
	})
	app.Logger.Fatal(app.Start())
}
