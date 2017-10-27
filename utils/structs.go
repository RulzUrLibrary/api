package utils

import (
	"encoding/gob"
	"fmt"
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
	Id    int64  `json:"-"`
	Name  string `json:"name"`
	Uuid  string `json:"uuid,omitempty"`
	User  string `json:"user,omitempty"`
	Books *Books `json:"books,omitempty"`
}

type Wishlists []Wishlist

func init() {
	gob.Register(&User{})
	gob.Register(Flash{})
}
