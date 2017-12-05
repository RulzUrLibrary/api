package utils

import (
	"encoding/gob"
	"fmt"
	"strconv"
	"strings"
)

type Dict map[string]interface{}
type flash int

const (
	FlashSuccess flash = iota
	FlashWarning
	FlashError
)

type Flash struct {
	Type flash
	Msg  string
}

// this should not be here
func (f Flash) Class() string {
	switch f.Type {
	case FlashSuccess:
		return "text-success"
	case FlashWarning:
		return "text-warning"
	case FlashError:
		return "text-error"
	}
	return ""
}

func (f Flash) Icon() string {
	switch f.Type {
	case FlashSuccess:
		return "icon-check"
	case FlashWarning:
		return "icon-time"
	case FlashError:
		return "icon-stop"
	}
	return ""
}

type User struct {
	Id    int
	Email string
}

type Book struct {
	Isbn        string     `json:"isbn"`
	Thumbnail   string     `json:"thumbnail,omitempty"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	Price       float32    `json:"price,omitempty"`
	Number      int        `json:"number,omitempty"`
	Serie       string     `json:"serie,omitempty"`
	Owned       *bool      `json:"owned,omitempty"` // tricking golang json encoding
	Authors     *Authors   `json:"authors,omitempty"`
	Wishlists   *Wishlists `json:"wishlists,omitempty"`
	Notations   *Notations `json:"notations,omitempty"`
}

func (b Book) TitleDisplay() string {
	if b.Serie == "" {
		return b.Title
	} else if b.Title == "" {
		return fmt.Sprintf("%s - %d", b.Serie, b.Number)
	} else {
		return fmt.Sprintf("%s - %d: %s", b.Serie, b.Number, b.Title)
	}
}

func (b Book) TitleDisplaySC() string {
	if b.Serie == "" {
		return b.Title
	} else {
		return fmt.Sprintf("%s tome %d", b.Serie, b.Number)
	}
}

func (b Book) InCollection() bool {
	return b.Owned != nil && *b.Owned
}

type Books []Book

func (b Books) owned() (nb int) {
	for _, book := range b {
		if book.InCollection() {
			nb++
		}
	}
	return
}

func (b Books) Owned() string {
	return fmt.Sprintf("%02d / %02d", b.owned(), len(b))
}

func (b Books) Ratio() float64 {
	return float64(b.owned()) / float64(len(b))
}

type Serie struct {
	Id          int64    `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Thumbnail   string   `json:"thumbnail,omitempty"`
	Title       string   `json:"title,omitempty"`
	Isbn        string   `json:"isbn,omitempty"`
	Owned       *bool    `json:"owned,omitempty"` // tricking golang json encoding
	Authors     *Authors `json:"authors,omitempty"`
	Volumes     *Books   `json:"volumes,omitempty"`
}

func (s Serie) Thumb() string {
	if s.Volumes == nil || len(*s.Volumes) == 0 {
		return "/thumbs/" + s.Isbn + ".jpg"
	}
	return "/thumbs/" + (*s.Volumes)[0].Isbn + ".jpg"
}

type Series []Serie

type Author struct {
	Id   uint64 `json:"-"`
	Name string `json:"name"`
}

type Authors []Author

func (a Authors) String() string {
	var names []string
	for _, author := range a {
		names = append(names, author.Name)
	}
	return strings.Join(names, ", ")
}

type Wishlist struct {
	Id      int64  `json:"-"`
	Checked bool   `json:"-"`
	Name    string `json:"name"`
	Uuid    string `json:"uuid,omitempty"`
	User    string `json:"user,omitempty"`
	Books   *Books `json:"books,omitempty"`
}

func (w Wishlist) EmptyBooks() Wishlist {
	if w.Books == nil {
		w.Books = &Books{}
	}
	return w
}

type Wishlists []Wishlist

func (w Wishlists) EmptyBooks() Wishlists {
	wishlists := make(Wishlists, len(w))
	for i, wishlist := range w {
		wishlists[i] = wishlist.EmptyBooks()
	}
	return wishlists
}

func (w Wishlists) Populate(book Book) *Wishlists {
	if book.Wishlists == nil || len(*book.Wishlists) == 0 {
		return &w
	}

	wishlists := make(map[string]*Wishlist, len(w))
	for i, wishlist := range w {
		wishlists[wishlist.Uuid] = &w[i]
	}
	for _, wishlist := range *book.Wishlists {
		wishlists[wishlist.Uuid].Checked = true
	}

	return &w
}

type Notation struct {
	Provider string  `json:"provider"`
	Note     float32 `json:"note"`
	Max      int     `json:"max"`
	Link     string  `json:"link"`
}

func (n Notation) DisplayNote() string {
	return fmt.Sprintf("%s / %d", strconv.FormatFloat(float64(n.Note), 'f', -1, 32), n.Max)
}

type Notations []Notation

func init() {
	gob.Register(&User{})
	gob.Register(Flash{})
}
