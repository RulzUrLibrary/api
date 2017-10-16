package app

import (
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

func WEBBookList(c *Context) (err error) {
	var series *db.Series
	var user = c.Get("user").(*utils.User)
	var query = struct {
		Isbn string `query:"isbn"`
		Pagination
	}{"", NewPagination()}

	if err = c.Bind(&query); err != nil {
		return
	}

	if err = c.Validate(&query); err != nil {
		return
	}

	if query.Isbn != "" {
		_, err = c.DB.BookPut([]string{query.Isbn}, user.Id)
		if err != nil {
			return
		}
		err = c.Flashes(utils.Flash{utils.FlashSuccess, "book added to collection!"})
		if err != nil {
			return
		}
		return c.Redirect(http.StatusSeeOther, c.Echo().Reverse("book", query.Isbn))
	}

	series, query.Count, err = c.DB.SerieListU(query.Limit(), query.Offset(), user.Id)
	if err != nil {
		return
	}

	return c.Render(http.StatusOK, "books.html", map[string]interface{}{
		"series":     series.ToStructs(true),
		"pagination": query.Pagination,
	})
}

func WEBBookGet(c *Context) error {
	book, err := BookGet(c)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "book.html", map[string]interface{}{"book": book})
}

func WEBSerieGet(c *Context) error {
	isbn := c.QueryParam("isbn")

	if isbn != "" {
		var user = c.Get("user").(*utils.User)
		if _, err := c.DB.BookPut([]string{isbn}, user.Id); err != nil {
			return err
		}
	}

	serie, err := SerieGet(c)
	if err != nil {
		return err
	}

	return c.Render(http.StatusOK, "serie.html", map[string]interface{}{"serie": serie})
}
