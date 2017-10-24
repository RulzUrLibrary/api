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
	if isbn := c.FormValue("isbn"); isbn != "" {
		var user = c.Get("user").(*utils.User)
		if _, err := c.Echo().Database.BookPut(user.Id, isbn); err != nil {
			return err
		}
	}
	return WEBSerieGet(c)
}
