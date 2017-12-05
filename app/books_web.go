package app

import (
	"github.com/RulzUrLibrary/api/ext/db"
	"github.com/RulzUrLibrary/api/utils"
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
	user, ok := c.Get("user").(*utils.User)
	if !ok {
		user = &utils.User{Id: 0}
	}
	if book, err := BookGet(c); err != nil {
		return err
	} else if wishlists, err := c.App.Database.WishlistsN(user.Id); err != nil {
		return err
	} else {
		book.Wishlists = wishlists.ToStructs(true).Populate(book)
		return c.Render(http.StatusOK, "book.html", dict{"book": book})
	}
}

func WEBBookPost(c *Context) error {
	var success, failure string

	req := struct {
		Wishlists []string `form:"wishlists"`
		Toggle    bool     `form:"toggle"`
	}{}
	user := c.Get("user").(*utils.User)
	isbn := c.Param("isbn")
	err := c.Bind(&req)

	if err != nil {
		return err
	}
	if req.Toggle {
		success = "collection_update_success"
		failure = "collection_update_failure"
		err = c.App.Database.CollectionToggle(user.Id, isbn)
	} else {
		success = "wishlist_update_success"
		failure = "wishlist_update_failure"
		err = c.App.Database.WishlistUpdate(user.Id, isbn, req.Wishlists...)
	}

	if err != nil {
		err = c.Flashes(utils.Flash{utils.FlashWarning, failure})
	} else {
		err = c.Flashes(utils.Flash{utils.FlashSuccess, success})
	}
	if err != nil {
		return err
	}
	return WEBBookGet(c)
}
