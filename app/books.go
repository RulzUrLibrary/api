package app

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

func BookGet(c *Context) (_ *utils.Book, err error) {
	var book *db.Book
	isbn := c.Param("isbn")
	user, ok := c.Get("user").(*utils.User)
	if ok {
		book, err = c.App.Database.BookGetU(isbn, user.Id)
	} else {
		book, err = c.App.Database.BookGet(isbn)
	}
	if err == utils.ErrNotFound {
		return nil, echo.NewHTTPError(http.StatusNotFound, "book "+isbn+" not found")
	}
	return book.ToStructs(false), err
}

func BookPost(c *Context, isbn string) (book utils.Book, ok bool, err error) {
	fn := func(book *db.Book, err error) (utils.Book, error) {
		if err != nil {
			return utils.Book{}, err
		}
		return *book.ToStructs(false), err
	}
	db := c.App.Database

	isbn = utils.SanitizeIsbn(isbn)
	if len(isbn) == 0 {
		err = echo.NewHTTPError(http.StatusBadRequest, "you provided an empty isbn")
		return
	}
	// check if book exists in database and return it if so
	if book, err = fn(db.BookGet(isbn)); err == nil {
		return book, true, nil
	} else if err != utils.ErrNotFound {
		return
	}
	// request additional informations
	switch book, err = c.App.Scrapper.Amazon(isbn); err {
	case nil:
		err = db.BookSave(&book)
		return
	case utils.ErrCaptcha:
		if _, err = db.CaptchaAdd(isbn); err == nil {
			err = echo.NewHTTPError(http.StatusAccepted,
				"request correctly received but unable to be processed currently.")
		}
		return
	case utils.ErrNoProduct:
		err = echo.NewHTTPError(http.StatusNotFound,
			fmt.Sprintf("product with isbn: '%s' not found", isbn))
		return
	}
	return
}

func BookList(c *Context, limit, offset int) ([]*utils.Book, int, error) {
	user, ok := c.Get("user").(*utils.User)
	if ok {
		return c.App.Database.BookListU(limit, offset, user.Id)
	} else {
		return c.App.Database.BookList(limit, offset)
	}
}

func change(c *Context, fn func(int, ...string) (int, error)) (int, error) {
	var user = c.Get("user").(*utils.User)
	var books struct {
		Isbns []string `json:"isbns"`
	}

	if err := c.Bind(&books); err != nil {
		return 0, err
	}
	if len(books.Isbns) == 0 {
		return 0, nil
	}
	c.App.Logger.Debug(books)
	return fn(user.Id, books.Isbns...)
}
