package app

import (
	"net/http"
)

func WEBBookGet(c *Context) error {
	book, err := BookGet(c)
	if err == nil {
		c.Render(http.StatusOK, "book.html", map[string]interface{}{"book": book})
	}
	return err
}
