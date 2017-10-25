package app

import (
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
	type st struct {
		success string
		failure string
		fn      func(int, ...string) (int, error)
	}
	s := map[string]map[string]st{
		"wishlist": map[string]st{
			"del": st{"book_wishlist_removed", "book_wishlist_already_removed", c.App.Database.BookDelete},
			"add": st{"book_wishlist_added", "book_wishlist_already_added", c.App.Database.WishlistPut},
		},
		"collection": map[string]st{
			"del": st{"book_collection_removed", "book_collection_already_removed", c.App.Database.BookDelete},
			"add": st{"book_collection_added", "book_collection_already_added", c.App.Database.BookPut},
		},
	}[c.FormValue("tag")][c.FormValue("action")]

	count, err := s.fn(c.Get("user").(*utils.User).Id, c.Param("isbn"))
	if err != nil {
		return err
	}

	if count == 0 {
		err = c.Flashes(utils.Flash{utils.FlashWarning, s.failure})
	} else {
		err = c.Flashes(utils.Flash{utils.FlashSuccess, s.success})
	}
	if err != nil {
		return err
	}
	return WEBBookGet(c)
}
