package app

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/paul-bismuth/library/utils"
	"net/http"
	"strconv"
)

type Serie struct {
	Name    string        `json:"name"`
	Authors utils.Authors `json:"authors"`
	Books   interface{}   `json:"volumes,omitempty"`
}

func SerieGet(c *Context) (interface{}, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return nil, echo.NewHTTPError(
			http.StatusBadRequest, "serie 'id' must be an integer",
		)
	}
	user, ok := c.Get("user").(utils.User)
	serie, books, err := c.DB.SerieGet(id, user.Id)
	if err != nil && err == utils.ErrNotFound {
		return nil, echo.NewHTTPError(
			http.StatusNotFound, fmt.Sprintf("serie %d not found", id),
		)
	}
	res := &Serie{Name: serie.Name, Authors: serie.Authors}
	if ok {
		res.Books = books.GetsS()
	} else {
		res.Books = books.Gets()
	}
	return res, nil
}
