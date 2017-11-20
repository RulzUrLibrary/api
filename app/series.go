package app

import (
	"github.com/labstack/echo"
	"github.com/RulzUrLibrary/api/ext/db"
	"github.com/RulzUrLibrary/api/utils"
	"net/http"
	"strconv"
)

func SerieGet(c *Context) (books db.Books, err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return db.Books{}, echo.NewHTTPError(
			http.StatusBadRequest, "serie 'id' must be an integer",
		)
	}

	user, ok := c.Get("user").(*utils.User)
	if ok {
		books, err = c.App.Database.SerieGetU(id, user.Id)
	} else {
		books, err = c.App.Database.SerieGet(id)
	}
	if err != nil {
		return
	}
	if len(books.Books) == 0 {
		return db.Books{}, echo.NewHTTPError(
			http.StatusNotFound, "serie "+c.Param("id")+" not found",
		)
	}
	return books, nil
}

func SerieList(c *Context, limit, offset int) (db.Books, int64, error) {
	user, ok := c.Get("user").(*utils.User)

	if ok {
		return c.App.Database.SerieListU(limit, offset, user.Id)
	} else {
		return c.App.Database.SerieList(limit, offset)
	}
}
