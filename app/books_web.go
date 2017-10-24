package app

import (
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

func WEBBookList(c *Context) (err error) {
	var series *db.Series
	var user = c.Get("user").(*utils.User)
	var query = NewPagination()

	if err = c.Bind(&query); err != nil {
		return
	}

	if err = c.Validate(&query); err != nil {
		return
	}

	series, query.Count, err = c.Echo().Database.SerieListU(query.Limit(),
		query.Offset(), user.Id)
	if err != nil {
		return
	}

	return c.Render(http.StatusOK, "books.html",
		dict{"series": series.ToStructs(true), "pagination": query})
}

func WEBBookGet(c *Context) error {
	if book, err := BookGet(c); err != nil {
		return err
	} else {
		return c.Render(http.StatusOK, "book.html", dict{"book": book})
	}
}

func WEBBookPost(c *Context) (err error) {
	var count int
	var user = c.Get("user").(*utils.User)

	if isbn := c.Param("isbn"); isbn != "" {
		count, err = c.Echo().Database.BookPut(user.Id, isbn)
		if err != nil {
			return
		}
		if count == 0 {
			err = c.Flashes(utils.Flash{utils.FlashWarning, "book already in collection!"})
		} else {
			err = c.Flashes(utils.Flash{utils.FlashSuccess, "book added to collection!"})
		}
		if err != nil {
			return
		}
	}
	return WEBBookGet(c)
}
