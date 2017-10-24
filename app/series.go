package app

import (
	"github.com/labstack/echo"
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
	"strconv"
)

func SerieGet(c *Context) (*utils.Serie, error) {
	var serie *db.Serie

	id, err := strconv.Atoi(c.Param("id"))
	user, ok := c.Get("user").(*utils.User)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusBadRequest, "serie 'id' must be an integer",
		)
	}
	if ok {
		serie, err = c.Echo().Database.SerieGetU(id, user.Id)
	} else {
		serie, err = c.Echo().Database.SerieGet(id)
	}
	switch err {
	case nil:
		return serie.ToStructs(false), nil
	case utils.ErrNotFound:
		return nil, echo.NewHTTPError(http.StatusNotFound, "serie "+c.Param("id")+" not found")
	}
	return nil, err
}

func SerieList(c *Context, limit, offset int) (*db.Series, int, error) {
	user, ok := c.Get("user").(*utils.User)

	if ok {
		return c.Echo().Database.SerieListU(limit, offset, user.Id)
	} else {
		return c.Echo().Database.SerieList(limit, offset)
	}
}
