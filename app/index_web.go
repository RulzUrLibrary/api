package app

import (
	"github.com/paul-bismuth/library/utils"
	"net/http"
)

type WEBSearch struct {
	Pattern string `query:"search"`
}

func WEBIndex(c *Context) (err error) {
	var search WEBSearch
	var books []*utils.Book

	if err = c.Bind(&search); err != nil {
		return
	}
	if search.Pattern == "" {
		books, _, err = c.DB.BookList(10, 0)
	} else {
		books, err = c.DB.BookSearch(search.Pattern, 10, 0)
	}
	if err != nil {
		return
	}
	c.Logger().Error(books)
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"books": books,
	})
}
