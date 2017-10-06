package app

import (
	"github.com/labstack/echo"
	"net/http"
)

func ContentType(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()

		if req.Method == echo.GET {
			return next(c)
		}
		ct := req.Header.Get(echo.HeaderContentType)
		if ct != echo.MIMEApplicationJSON {
			return echo.NewHTTPError(
				http.StatusBadRequest, "API only support application/json content type",
			)
		}
		return next(c)
	}
}
