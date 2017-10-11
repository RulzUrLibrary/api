package app

import (
	"github.com/paul-bismuth/library/utils"
	"net/http"
	"strings"
)

type WEBSearch struct {
	Pattern string `query:"search"`
}

func WEBIndex(c *Context) (err error) {
	var search WEBSearch
	var books []*utils.Book
	var pattern string

	if err = c.Bind(&search); err != nil {
		return
	}
	pattern = strings.TrimSpace(search.Pattern)
	if utils.IsIsbn10(pattern) || utils.IsIsbn13(pattern) {
		book, _, err := BookPost(c, pattern)
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "book.html", map[string]interface{}{
			"book": book,
		})
	} else if pattern == "" {
		books, _, err = c.DB.BookList(10, 0)
	} else {
		books, err = c.DB.BookSearch(pattern, 10, 0)
	}
	if err != nil {
		return
	}
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"books": books,
	})
}
