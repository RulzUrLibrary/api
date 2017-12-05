package app

import (
	"github.com/RulzUrLibrary/api/ext/db"
	"github.com/RulzUrLibrary/api/utils"
	"net/http"
)

func APIWishlist(c *Context) (err error) {
	var wishlists db.Wishlists
	var meta = NewMeta()
	var user = c.Get("user").(*utils.User).Id

	if err = c.Bind(&meta); err != nil {
		return
	}
	if err = c.Validate(&meta); err != nil {
		return
	}

	wishlists, meta.Count, err = c.App.Database.Wishlists(meta.Limit, meta.Offset, user)
	if err != nil {
		return
	}
	return c.JSON(http.StatusOK, dict{"_meta": meta, "wishlists": wishlists.ToStructs(true).EmptyBooks()})
}

func APIWishlistPost(c *Context) error {
	isbn := c.Param("isbn")
	user := c.Get("user").(*utils.User)
	request := struct {
		Wishlists []string `json:"wishlists"`
	}{}

	if err := c.Bind(&request); err != nil {
		return err
	}

	if err := c.App.Database.WishlistUpdate(user.Id, isbn, request.Wishlists...); err != nil {
		return err
	}

	c.Response().WriteHeader(http.StatusNoContent)
	return nil
}
