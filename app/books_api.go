package app

import (
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
	s := newSearch()

	if err := c.Bind(&s); err != nil {
		return err
	}
	if s.Pattern == "" {
		books, count, err := BookList(c, int(s.Limit), int(s.Offset))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"_meta": Meta{int(s.Limit), int(s.Offset), count}, "books": books,
		})
	} else {
		books, err := c.DB.BookSearch(s.Pattern, int(s.Limit), int(s.Offset))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]interface{}{
			"_meta": Meta{int(s.Limit), int(s.Offset), -1}, "books": books,
		})
	}
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
