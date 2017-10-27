package app

import (
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

func APIBookGet(c *Context) error {
	book, err := BookGet(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, book)
}

func APIBookPost(c *Context) error {
	b := struct {
		Isbn string `json:"isbn"`
	}{}
	err := c.Bind(&b)
	if err != nil {
		return err
	}
	i, ok, err := BookPost(c, b.Isbn)
	if err != nil {
		return err
	}
	if ok {
		return c.JSON(http.StatusOK, i)
	} else {
		return c.JSON(http.StatusCreated, i)
	}
}

func APIBookList(c *Context) (err error) {
	var books db.Books
	var s = struct {
		Pattern string `query:"search"`
		Meta
	}{"", NewMeta()}

	if err = c.Bind(&s); err != nil {
		return
	}
	if err = c.Validate(&s); err != nil {
		return
	}
	if s.Pattern == "" {
		user, ok := c.Get("user").(*utils.User)
		if ok {
			books, s.Count, err = c.App.Database.BookListU(s.Limit, s.Offset, user.Id)
		} else {
			books, s.Count, err = c.App.Database.BookList(s.Limit, s.Offset)
		}
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, dict{"_meta": s.Meta, "books": books.ToStructs(true)})
	} else {
		books, err = c.App.Database.BookSearch(s.Pattern, s.Limit, s.Offset)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, dict{"_meta": s.Meta, "books": books.ToStructs(true)})
	}
}

func APIBookPut(c *Context) error {
	count, err := change(c, c.App.Database.BookPut)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, struct {
		Count int64 `json:"added"`
	}{count})
}

func APIBookDelete(c *Context) (err error) {
	count, err := change(c, c.App.Database.BookDelete)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, struct {
		Count int64 `json:"deleted"`
	}{count})
}
