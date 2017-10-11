package app

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paul-bismuth/library/ext/db"
	"github.com/paul-bismuth/library/utils"
	"net/http"
)

func BookGet(c *Context) (_ *utils.Book, err error) {
	var book *db.Book
	isbn := c.Param("isbn")
	user, ok := c.Get("user").(*utils.User)
	if ok {
		book, err = c.DB.BookGetU(isbn, user.Id)
	} else {
		book, err = c.DB.BookGet(isbn)
	}
	if err == utils.ErrNotFound {
		return nil, echo.NewHTTPError(http.StatusNotFound, "book "+isbn+" not found")
	}
	return book.ToStructs(false), err
}

func BookPost(c *Context, isbn string) (_ interface{}, ok bool, err error) {
	isbn = utils.SanitizeIsbn(isbn)
	if len(isbn) == 0 {
		return nil, ok, echo.NewHTTPError(
			http.StatusBadRequest, "you provided an empty isbn",
		)
	}
	// check if book exists in database and return it if so
	if book, err := c.DB.BookGet(isbn); err == nil {
		return book.ToStructs(false), true, nil
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

func BookList(c *Context, limit, offset int) (_ map[string]interface{}, err error) {
	var books utils.Books
	var count int

	user, ok := c.Get("user").(*utils.User)
	if ok {
		books, count, err = c.DB.BookListU(limit, offset, user.Id)
	} else {
		books, count, err = c.DB.BookList(limit, offset)
	}
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"_meta": Meta{limit, offset, count}, "books": books,
	}, nil
}

func BookSearch(c *Context, pattern string, limit, offset int) (map[string]interface{}, error) {
	books, err := c.DB.BookSearch(pattern, limit, offset)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"_meta": Meta{limit, offset, -1}, "books": books,
	}, nil
}

func change(c *Context, fn func([]string, int) (int, error)) (int, error) {
	var user = c.Get("user").(*utils.User)
	var books struct {
		Isbns []string `json:"isbns" query:"isbn"`
	}

	if err := c.Bind(&books); err != nil {
		return 0, err
	}
	if len(books.Isbns) == 0 {
		return 0, nil
	}
	c.Logger.Debug(books)
	return fn(books.Isbns, user.Id)
}
