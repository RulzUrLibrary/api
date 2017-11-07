package app

import (
	"github.com/labstack/echo"
	"github.com/rulzurlibrary/api/ext/db"
	"github.com/rulzurlibrary/api/ext/validator"
	"github.com/rulzurlibrary/api/utils"
	"net/http"
)

func WEBWishlist(c *Context) (err error) {
	var wishlists db.Wishlists
	var query = NewPagination()
	var user = c.Get("user").(*utils.User).Id

	if err = c.Bind(&query); err != nil {
		return
	}
	if err = c.Validate(&query); err != nil {
		return
	}

	wishlists, query.Count, err = c.App.Database.Wishlists(query.Limit(), query.Offset(), user)
	if err != nil {
		return
	}
	return c.Render(http.StatusOK, "wishlists.html", dict{
		"wishlists": wishlists.ToStructs(true), "pagination": query,
		"links": wishlists.AbsLinks(c.ReverseAbs),
	})
}

func WEBTag(c *Context) (err error) {
	//var books []*utils.Book
	form := struct {
		Name string `form:"name" validate:"required"`
	}{}
	render := func(code int, errs dict) error {
		return c.Render(code, "tags.html", dict{"form": form, "errors": errs})
	}

	if c.Request().Method == http.MethodPost {
		if err = c.Bind(&form); err != nil {
			return
		}
		if err = c.Validate(&form); err != nil {
			return render(http.StatusBadRequest, validator.Dump(err, map[string]dict{
				"name": dict{"required": utils.WISHLIST_NAME_REQUIRED},
			}))
		}

		err = c.App.Database.WishPost(form.Name, c.Get("user").(*utils.User).Id)
		switch err {
		case nil:
			err = c.Flashes(utils.Flash{utils.FlashSuccess, utils.WISHLIST_CREATED})
			if err != nil {
				return err
			}
			return c.Redirect(http.StatusSeeOther, c.Echo().Reverse("wishlists"))
		case utils.ErrExists:
			return render(http.StatusBadRequest, dict{"name": utils.WISHLIST_NAME_ALREADY_TAKEN})
		default:
			return
		}
	}
	return render(http.StatusOK, dict{})
}

func WEBWishListGet(c *Context) (err error) {
	var wishlists db.Wishlists
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

	wishlists, query.Count, err = c.App.Database.Wishlist(
		query.Limit(), query.Offset(), uuid,
	)
	if err != nil {
		return
	}
	return c.Render(http.StatusOK, "wishlist.html", dict{
		"wishlist": wishlists.ToStructs(false)[0], "pagination": query,
	})
}

func WEBWishListPost(c *Context) (err error) {
	var uuid = c.Param("id")
	var user = c.Get("user").(*utils.User)

	count, err := c.App.Database.WishlistDelete(user.Id, uuid)
	if err != nil {
		return err
	}
	flash := utils.Flash{utils.FlashWarning, utils.WISHLIST_DELETE_WARNING}
	if count > 0 {
		flash = utils.Flash{utils.FlashSuccess, utils.WISHLIST_DELETE_SUCCESS}
	}
	if err := c.Flashes(flash); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, c.Reverse("wishlists"))
}

func WEBWishlistAdd(c *Context) error {
	isbn := c.Param("isbn")
	user := c.Get("user").(*utils.User)

	book, err := c.App.Database.BookGet(isbn)
	if err != nil {
		return err
	}
	wishlists, err := c.App.Database.WishlistsN(user.Id)
	switch err {
	case nil:
	case utils.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, "book "+isbn+" not found")
	default:
		return err
	}

	return c.Render(http.StatusOK, "add.html",
		dict{"book": book.ToStructs(false), "wishlists": wishlists.ToStructs(true)})
}

func WEBWishlistPost(c *Context) error {
	isbn := c.Param("isbn")
	user := c.Get("user").(*utils.User)
	wishlists := []string{""}

	if _wishlists, ok := c.Request().Form["wishlists"]; ok {
		wishlists = _wishlists
	}

	count, err := c.App.Database.WishlistPut(user.Id, isbn, wishlists...)
	if err != nil {
		return err
	}
	flash := utils.Flash{utils.FlashWarning, utils.WISHLIST_ADD_WARNING}
	if count > 0 {
		flash = utils.Flash{utils.FlashSuccess, utils.WISHLIST_ADD_SUCCESS}
	}
	if err := c.Flashes(flash); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, c.Reverse("book", isbn))
}
