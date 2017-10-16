package app

import (
	"github.com/rulzurlibrary/api/ext/db"
	"net/http"
)

func APISerieGet(c *Context) error {
	serie, err := SerieGet(c)
	if err == nil {
		c.JSON(http.StatusOK, serie)
	}
	return err
}

func APISerieList(c *Context) (err error) {
	var series *db.Series
	var meta = NewMeta()

	if err = c.Bind(&meta); err != nil {
		return
	}
	if err = c.Validate(&meta); err != nil {
		return
	}
	series, meta.Count, err = SerieList(c, meta.Limit, meta.Offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"_meta": meta, "series": series.ToStructs(true),
	})
}
