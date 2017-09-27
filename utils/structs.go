package utils

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
)

type User struct {
	Id   int
	Name string
}

type Book struct {
	Id          int     `json:"-"`
	Isbn        string  `json:"isbn"`
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description"`
	Thumbnail   string  `json:"thumbnail"`
	Price       float32 `json:"price,omitempty"`
	Authors     Authors `json:"authors"`
	Number      int     `json:"number,omitempty"`
	Serie       string  `json:"serie,omitempty"`
}

type BookScoped struct {
	Book
	Owned bool `json:"owned"`
}

func (b Book) String() string {
	marshaled, err := json.Marshal(b)
	if err != nil {
		return fmt.Sprintf("error marshaling book: %s", err)
	}
	return string(marshaled)
}

func (b Book) TitleDisplay() string {
	return TitleDisplay(b.Title, b.Serie, b.Number)
}

func (b Book) ToArgs() []interface{} {
	args := []interface{}{
		b.Isbn, b.Title, b.Description, b.Price, b.Number, nil,
	}

	if b.Title == "" {
		args[1] = nil
	}
	if b.Number == 0 {
		args[4] = nil
	}
	return args
}

func TitleDisplay(title string, serie string, number int) string {
	if serie == "" {
		return title
	} else if title == "" {
		return fmt.Sprintf("%s - %d", serie, number)
	} else {
		return fmt.Sprintf("%s - %d: %s", serie, number, title)
	}
}

type Volume struct {
	Number int    `json:"number"`
	Isbn   string `json:"isbn"`
	Owned  bool   `json:"owned"`
}

type Volumes []Volume

func (v Volumes) owned() (nb int) {
	for _, volume := range v {
		if volume.Owned {
			nb++
		}
	}
	return
}

func (v Volumes) String() string {
	return fmt.Sprintf("%02d / %02d", v.owned(), len(v))
}

func (v Volumes) Ratio() float64 {
	return float64(v.owned()) / float64(len(v))
}

type Serie struct {
	Id          int     `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Title       string  `json:"title,omitempty"`
	Isbn        string  `json:"isbn,omitempty"`
	Authors     Authors `json:"authors"`
	Volumes     Volumes `json:"volumes,omitempty"`
}

func (s Serie) IsSerie() bool {
	return s.Isbn == ""
}

func (s Serie) Thumb() string {
	isbn := s.Isbn

	if s.IsSerie() {
		isbn = s.Volumes[0].Isbn
	}

	return "/thumbs/" + isbn + ".jpg"
}

type Author struct {
	Id   int    `json:"-"`
	Name string `json:"name"`
}

type Authors []*Author

func (a Authors) String() string {
	var names []string
	for _, author := range a {
		names = append(names, author.Name)
	}
	return strings.Join(names, ", ")
}

func init() {
	gob.Register(&User{})
}
