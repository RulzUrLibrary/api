package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/rulzurlibrary/api/utils"
)

type Wishlist struct {
	id   sql.NullInt64
	name sql.NullString
	uuid sql.NullString
	user sql.NullString
	book Book
}

func (w *Wishlist) Scan(src interface{}) (err error) {
	var elems [][]byte

	if elems, err = parseRow(src.([]byte), []byte{','}); err != nil {
		return
	}

	if len(elems) != 2 {
		return fmt.Errorf("element is not a valid wishlist")
	}

	w.name.String = string(elems[0])
	w.uuid.String = string(elems[1])
	return
}

type Wishlists struct {
	Wishlists []*Wishlist
	Valid     bool
}

func (w *Wishlists) InsertWishlist(fn func(*Wishlist) list) func() list {
	return func() list {
		var wishlist Wishlist
		w.Wishlists = append(w.Wishlists, &wishlist)
		return fn(&wishlist)
	}
}

// You are responsible for wishlists ids to be consecutives
func (w *Wishlists) ToStructs(partial bool) (wishlists utils.Wishlists) {
	var last utils.Wishlist

	for _, wishlist := range w.Wishlists {
		if last.Id != 0 && wishlist.id.Int64 == last.Id {
			*last.Books = append(*last.Books, wishlist.book.ToStructs(partial))
		} else {
			last = utils.Wishlist{
				wishlist.id.Int64, wishlist.name.String, wishlist.uuid.String,
				wishlist.user.String, nil}
			if wishlist.book.id.Valid {
				last.Books = &utils.Books{wishlist.book.ToStructs(partial)}
			}
			wishlists = append(wishlists, last)
		}
	}
	return
}

func (w *Wishlists) ToWishlists(partial bool) *utils.Wishlists {
	if !w.Valid {
		return nil
	}
	wishlists := w.ToStructs(partial)
	if wishlists == nil {
		return &utils.Wishlists{}
	}
	return &wishlists
}

func (w *Wishlists) AbsLinks(fn func(string, ...interface{}) string) map[string]string {
	links := make(map[string]string)
	for _, wishlist := range w.Wishlists {
		links[wishlist.uuid.String] = fn("wishlist", wishlist.uuid.String)
	}
	return links
}

func (w *Wishlists) Scan(src interface{}) error {
	w.Valid = true
	if bytes.Equal(src.([]byte), []byte(`{"(,)"}`)) {
		return nil
	}
	return pq.Array(&w.Wishlists).Scan(src)
}
