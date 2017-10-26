package app

import (
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

func WEBWishList(c *Context) (err error) {
	var books []*utils.Book
	var query = NewPagination()

	if err = c.Bind(&query); err != nil {
		return
	}
	if err = c.Validate(&query); err != nil {
		return
	}

	books, query.Count, err = c.App.Database.WishList(query.Limit(),
		query.Offset(), c.Get("user").(*utils.User).Id)

	if err != nil {
		return
	}

	return c.Render(http.StatusOK, "wishlist.html", dict{"books": books, "pagination": query})
}

func WEBWishListShare(c *Context) error {
	uuid, err := c.App.Database.WishListLink(c.Get("user").(*utils.User).Id)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "share.html", dict{"link": c.ReverseAbs("share", uuid)})
}

func WEBWishListGet(c *Context) (err error) {
	var books []*utils.Book
	var query = NewPagination()
	var uuid = c.Param("id")

	if err = dynamic(c, "wishlists", "uuid", uuid); err != nil {
		return
	}
	if err = c.Bind(&query); err != nil {
		return
	}
	if err = c.Validate(&query); err != nil {
		return
	}

	books, query.Count, err = c.App.Database.WishListU(query.Limit(), query.Offset(), uuid)
	if err != nil {
		return
	}

	return c.Render(http.StatusOK, "shared.html", dict{"books": books, "pagination": query})
}
