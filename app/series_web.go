package app

import (
	"github.com/RulzUrLibrary/api/utils"
	"net/http"
)

func WEBSerieGet(c *Context) error {
	if books, err := SerieGet(c); err != nil {
		return err
	} else {
		return c.Render(http.StatusOK, "serie.html",
			dict{"serie": books.ToSeries(false)[0]})
	}
}

func WEBSeriePost(c *Context) error {
	var user = c.Get("user").(*utils.User)

	switch c.FormValue("action") {
	case "collection":
		if _, err := c.App.Database.BookPut(user.Id, c.FormValue("isbn")); err != nil {
			return err
		}
	case "wishlist":
		if _, err := c.App.Database.WishlistPut(user.Id, c.FormValue("isbn")); err != nil {
			return err
		}
	}
	return WEBSerieGet(c)
}
