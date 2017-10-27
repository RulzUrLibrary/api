package app

import (
	"github.com/labstack/echo"
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

func WEBBookList(c *Context) (err error) {
	var series db.Books
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
		dict{"series": series.ToSeries(true), "pagination": query})
}

func WEBBookGet(c *Context) error {
	if book, err := BookGet(c); err != nil {
		return err
	} else {
		return c.Render(http.StatusOK, "book.html", dict{"book": book})
	}
}

func WEBBookPost(c *Context) (err error) {
	var success, failure string
	var count int64

	if wish := c.FormValue("wish"); wish != "" {
		success = "book_wishlist_removed"
		failure = "book_wishlist_already_removed"
		count, err = c.App.Database.WishDelete(c.Get("user").(*utils.User).Id,
			c.Param("isbn"), wish)
	} else if action := c.FormValue("action"); action == "del" {
		success = "book_collection_removed"
		failure = "book_collection_already_removed"
		count, err = c.App.Database.BookDelete(c.Get("user").(*utils.User).Id,
			c.Param("isbn"))
	} else if action == "add" {
		success = "book_collection_added"
		failure = "book_collection_already_added"
		count, err = c.App.Database.BookPut(c.Get("user").(*utils.User).Id,
			c.Param("isbn"))
	} else {
		err = echo.NewHTTPError(http.StatusBadRequest, nil)
	}
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
