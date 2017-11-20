package app

import (
	"github.com/RulzUrLibrary/api/ext/db"
	"github.com/RulzUrLibrary/api/utils"
	"net/http"
	"strings"
)

func WEBIndex(c *Context) (err error) {
	var books db.Books
	var pattern string
	var query = struct {
		Pattern string `query:"search"`
		Pagination
	}{"", NewPagination()}

	if err = c.Bind(&query); err != nil {
		return
	}
	if err = c.Validate(&query); err != nil {
		return
	}
	pattern = strings.TrimSpace(query.Pattern)
	if utils.IsIsbn10(pattern) || utils.IsIsbn13(pattern) {
		book, _, err := BookPost(c, pattern)
		if err != nil {
			return err
		}
		return c.Render(http.StatusOK, "book.html", dict{"book": book})
	} else if pattern == "" {
		books, query.Count, err = c.App.Database.BookList(query.Limit(), query.Offset())
	} else {
		books, err = c.App.Database.BookSearch(pattern, query.Limit(), query.Offset())
	}
	if err != nil {
		return
	}
	return c.Render(http.StatusOK, "index.html", dict{
		"books": books.ToStructs(false), "pagination": query.Pagination,
	})
}
