package app

import (
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

func WEBSerieGet(c *Context) error {
	if serie, err := SerieGet(c); err != nil {
		return err
	} else {
		return c.Render(http.StatusOK, "serie.html", dict{"serie": serie})
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
