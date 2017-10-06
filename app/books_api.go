package app

import (
	"github.com/paul-bismuth/library/utils"
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

type APISearch struct {
	Pattern string `query:"search"`
	Pagination
}

func newSearch() APISearch {
	return APISearch{Pagination: NewPagination()}
}

func APIBookList(c *Context) error {
	var res interface{}

	s := newSearch()
	err := c.Bind(&s)
	if err != nil {
		return err
	}
	if s.Pattern == "" {
		res, err = BookList(c, int(s.Limit), int(s.Offset))
	} else {
		res, err = BookSearch(c, s.Pattern, int(s.Limit), int(s.Offset))
	}
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, res)
}

func change(c *Context, fn func([]string, int) (int, error)) (int, error) {
	var user = c.Get("user").(utils.User)
	var books struct {
		Isbns []string `json:"isbns"`
	}

	if err := c.Bind(&books); err != nil {
		return 0, err
	}
	c.Logger.Debug(books)
	return fn(books.Isbns, user.Id)
}

func APIBookPut(c *Context) error {
	count, err := change(c, c.DB.BookPut)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, struct {
		Count int `json:"added"`
	}{count})
}

func APIBookDelete(c *Context) (err error) {
	count, err := change(c, c.DB.BookDelete)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, struct {
		Count int `json:"deleted"`
	}{count})
}
