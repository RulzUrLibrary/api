package app

import (
	"net/http"
)

func WEBBookList(c *Context) error {
	return c.Render(http.StatusOK, "books.html", map[string]interface{}{
		"books": nil,
	})
}

func WEBBookGet(c *Context) error {
	book, err := BookGet(c)
	if err == nil {
		c.Render(http.StatusOK, "book.html", map[string]interface{}{"book": book})
	}
	return err
}
