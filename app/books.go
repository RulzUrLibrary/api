package app

import (
	"fmt"
	"github.com/ixday/echo-hello/utils"
	"github.com/labstack/echo"
	"net/http"
)

func BookGet(c *Context) (interface{}, error) {
	isbn := c.Param("isbn")
	user, ok := c.Get("user").(utils.User)
	book, err := c.DB.BookGet(isbn, user.Id)
	if err != nil && err == utils.ErrNotFound {
		return nil, echo.NewHTTPError(http.StatusNotFound, "book "+isbn+" not found")
	}
	if ok {
		return book.ToStructsS(), nil
	} else {
		return book.ToStructs(), nil
	}
}
