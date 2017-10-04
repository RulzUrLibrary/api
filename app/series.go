package app

import (
	"github.com/labstack/echo"
	"github.com/paul-bismuth/library/ext/db"
	"github.com/paul-bismuth/library/utils"
	"net/http"
	"strconv"
)

func SerieGet(c *Context) (interface{}, error) {
	var serie *db.Serie

	id, err := strconv.Atoi(c.Param("id"))
	user, ok := c.Get("user").(utils.User)
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusBadRequest, "serie 'id' must be an integer",
		)
	}
	if ok {
		serie, err = c.DB.SerieGetU(id, user.Id)
		if err == nil {
			return serie.ToSerieGetScoped(), nil
		}
	} else {
		serie, err = c.DB.SerieGet(id)
		if err == nil {
			return serie.ToSerieGet(), nil
		}
	}
	if err == utils.ErrNotFound {
		return nil, echo.NewHTTPError(http.StatusNotFound, "serie "+c.Param("id")+" not found")
	}
	return nil, err
}

func SerieList(c *Context, limit, offset int) (_ interface{}, err error) {
	var series *db.Series

	user, ok := c.Get("user").(utils.User)

	res := struct {
		Meta   `json:"_meta"`
		Series interface{} `json:"series"`
	}{Meta{limit, offset, 0}, nil}
	if ok {
		series, res.Meta.Count, err = c.DB.SerieListU(limit, offset, user.Id)
		if err == nil {
			res.Series = series.ToSeriesScoped()
		}
	} else {
		series, res.Meta.Count, err = c.DB.SerieList(limit, offset)
		if err == nil {
			res.Series = series.ToSeries()
		}
	}
	return res, err
}
