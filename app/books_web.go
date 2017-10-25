package app

import (
	"github.com/labstack/echo"
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

func WEBBookList(c *Context) (err error) {
	var series *db.Series
	var query = NewPagination()

	if err = c.Bind(&query); err != nil {
		return
	}

	if err = c.Validate(&query); err != nil {
		return
	}

	series, query.Count, err = c.App.Database.SerieListU(query.Limit(),
		query.Offset(), c.Get("user").(*utils.User).Id)
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

func WEBBookPost(c *Context) error {
	var success, failure string
	var fn func(int, ...string) (int, error)

	switch c.FormValue("action") {
	case "del":
		success = "book_removed"
		failure = "book_already_removed"
		fn = c.App.Database.BookDelete
	case "add":
		success = "book_added"
		failure = "book_already_added"
		fn = c.App.Database.BookPut
	default:
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	count, err := fn(c.Get("user").(*utils.User).Id, c.Param("isbn"))
	if err != nil {
		return err
	}

	if count == 0 {
		err = c.Flashes(utils.Flash{utils.FlashWarning, failure})
	} else {
		err = c.Flashes(utils.Flash{utils.FlashSuccess, success})
	}
	if err != nil {
		return err
	}
	return WEBBookGet(c)
}
