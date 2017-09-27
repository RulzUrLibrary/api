package app

import (
	"github.com/ixday/echo-hello/ext/db"
	"github.com/ixday/echo-hello/utils"
	"github.com/labstack/echo"
	"net/http"
)

func BookGet(c *Context) (book *db.Book, err error) {
	isbn := c.Param("isbn")
	book, err = c.DB.BookGet(isbn)
	if err != nil && err == utils.ErrNotFound {
		err = echo.NewHTTPError(http.StatusNotFound, "book "+isbn+" not found")
	}
	return
}
