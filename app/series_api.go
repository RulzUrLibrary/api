package app

import (
	"net/http"
)

func APISerieGet(c *Context) error {
	serie, err := SerieGet(c)
	if err == nil {
		c.JSON(http.StatusOK, serie)
	}
	return err
}

func APISerieList(c *Context) error {
	p := NewPagination()
	err := c.Bind(&p)
	if err != nil {
		return err
	}
	series, err := SerieList(c, int(p.Limit), int(p.Offset))
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, series)
	return nil
}
