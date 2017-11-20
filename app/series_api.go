package app

import (
	"github.com/RulzUrLibrary/api/ext/db"
	"net/http"
)

func APISerieGet(c *Context) error {
	serie, err := SerieGet(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, serie.ToSeries(false)[0])
}

func APISerieList(c *Context) (err error) {
	var books db.Books
	var meta = NewMeta()

	if err = c.Bind(&meta); err != nil {
		return
	}
	if err = c.Validate(&meta); err != nil {
		return
	}
	books, meta.Count, err = SerieList(c, meta.Limit, meta.Offset)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dict{"_meta": meta, "series": books.ToSeries(true)})
}
