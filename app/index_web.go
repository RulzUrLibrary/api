package app

import (
	"github.com/paul-bismuth/library/utils"
	"net/http"
	"strings"
)

type WEBSearch struct {
	Pattern string `query:"search"`
	Pagination
}

func WEBIndex(c *Context) (err error) {
	var books []*utils.Book
	var pattern string
	var search = WEBSearch{"", NewPagination()}

	if err = c.Bind(&search); err != nil {
		return
	}
	if err = c.Validate(&search); err != nil {
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
		books, search.Count, err = c.DB.BookList(search.Limit(), search.Offset())
	} else {
		books, err = c.DB.BookSearch(pattern, search.Limit(), search.Offset())
	}
	if err != nil {
		return
	}
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"books":      books,
		"pagination": search.Pagination,
	})
}
