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
