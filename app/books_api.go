package app

import (
	"net/http"
)

func APIBookGet(c *Context) error {
	book, err := BookGet(c)
	if err == nil {
		c.JSON(http.StatusOK, book)
	}
	return err
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
		c.JSON(http.StatusOK, i)
	} else {
		c.JSON(http.StatusCreated, i)
	}
	return nil
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

	if err == nil {
		c.JSON(http.StatusOK, res)
	}
	return err
}
