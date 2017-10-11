package app

import (
	"fmt"
	"github.com/paul-bismuth/library/utils"
	"net/http"
)

func WEBBookList(c *Context) error {
	var isbns []string

	_, err := change(c, func(added []string, i int) (int, error) {
		isbns = added
		return c.DB.BookPut(isbns, i)
	})
	if err != nil {
		return nil
	}
	switch len(isbns) {
	case 0:
	case 1:
		flash := utils.Flash{utils.FlashSuccess, "book added to collection!"}
		if err := c.Flashes(flash); err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, c.Echo().Reverse("book", isbns[0]))
	default:
		flashes := []utils.Flash{}
		for _, isbn := range isbns {
			msg := fmt.Sprintf("book with isbn: %s, added to collection", isbn)
			flashes = append(flashes, utils.Flash{utils.FlashSuccess, msg})
		}
		if err := c.Flashes(flashes...); err != nil {
			return err
		}
	}
	return c.Render(http.StatusOK, "books.html", map[string]interface{}{
		"books": nil,
	})
}

func WEBBookGet(c *Context) error {
	book, err := BookGet(c)
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "book.html", map[string]interface{}{"book": book})
}
