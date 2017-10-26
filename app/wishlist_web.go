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
