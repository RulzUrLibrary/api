package app

import (
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
	var books []*utils.Book
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
		books, s.Count, err = BookList(c, s.Limit, s.Offset)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, dict{"_meta": s.Meta, "books": books})
	} else {
		books, err = c.App.Database.BookSearch(s.Pattern, s.Limit, s.Offset)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, dict{"_meta": s.Meta, "books": books})
	}
}

func APIBookPut(c *Context) error {
	count, err := change(c, c.App.Database.BookPut)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, struct {
		Count int `json:"added"`
	}{count})
}

func APIBookDelete(c *Context) (err error) {
	count, err := change(c, c.App.Database.BookDelete)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, struct {
		Count int `json:"deleted"`
	}{count})
}
