package app

import (
	"net/http"
)

func APIBookGet(c *Context) error {
	book, err := BookGet(c)
	if err == nil {
		c.JSON(http.StatusOK, book)
	}
	return err
}

func APIBookPost(c *Context) error {
	b := struct {
		Isbn string `json:"isbn"`
	}{}
	err := c.Bind(&b)
	if err != nil {
		return err
	}
	i, ok, err := BookPost(c, b.Isbn)
	if err != nil {
		return err
	}
	if ok {
		c.JSON(http.StatusOK, i)
	} else {
		c.JSON(http.StatusCreated, i)
	}
	return nil
}
