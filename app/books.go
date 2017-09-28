package app

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paul-bismuth/library/ext/db"
	"github.com/paul-bismuth/library/utils"
	"net/http"
)

func BookGet(c *Context) (interface{}, error) {
	isbn := c.Param("isbn")
	user, ok := c.Get("user").(utils.User)
	book, err := c.DB.BookGet(isbn, user.Id)
	if err != nil && err == utils.ErrNotFound {
		return nil, echo.NewHTTPError(http.StatusNotFound, "book "+isbn+" not found")
	}
	if ok {
		return book.ToStructsS(), nil
	} else {
		return book.ToStructs(), nil
	}
}

func BookPost(c *Context, isbn string) (_ interface{}, ok bool, err error) {
	isbn = utils.SanitizeIsbn(isbn)
	if len(isbn) == 0 {
		return nil, ok, echo.NewHTTPError(
			http.StatusBadRequest, "you provided an empty isbn",
		)
	}
	// check if book exists in database and return it if so
	if book, err := c.DB.BookGet(isbn, 0); err == nil {
		return book.ToStructs(), true, nil
	} else if err != utils.ErrNotFound {
		return nil, ok, err
	}
	// request additional informations
	var book utils.Book
	switch book, err = c.Scrapper.Amazon(isbn); err {
	case nil:
		err := c.DB.BookSave(&book)
		return book, ok, err
	case utils.ErrCaptcha:
		return nil, ok, echo.NewHTTPError(
			http.StatusAccepted,
			"request correctly received but unable to be processed currently.",
		)
	case utils.ErrNoProduct:
		return nil, ok, echo.NewHTTPError(
			http.StatusNotFound,
			fmt.Sprintf("product with isbn: '%s' not found", isbn),
		)
	}
	return nil, ok, err
}

func BookList(c *Context, limit, offset int) (interface{}, error) {
	count, err := c.DB.Count(db.CountBooks)
	if err != nil {
		return nil, err
	}

	books, err := c.DB.BookList(db.SelectBooks, limit, offset)
	if err != nil {
		return nil, err
	}
	return struct {
		Meta  `json:"_meta"`
		Books []*utils.Book `json:"books"`
	}{Meta{limit, offset, count}, books}, nil
}
