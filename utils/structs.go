package utils

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
)

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
	Title  string `json:"title,omitempty"`
	Number int    `json:"number"`
	Isbn   string `json:"isbn"`
}

type VolumeScoped struct {
	Volume
	Owned bool `json:"owned"`
}

type Volumes []*Volume
type VolumesScoped []*VolumeScoped

func (v VolumesScoped) owned() (nb int) {
	for _, volume := range v {
		if volume.Owned {
			nb++
		}
	}
	return
}

func (v VolumesScoped) String() string {
	return fmt.Sprintf("%02d / %02d", v.owned(), len(v))
}

func (v VolumesScoped) Ratio() float64 {
	return float64(v.owned()) / float64(len(v))
}

type Serie struct {
	Id          int     `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Thumbnail   string  `json:"thumbnail,omitempty"`
	Title       string  `json:"title,omitempty"`
	Isbn        string  `json:"isbn,omitempty"`
	Owned       *bool   `json:"owned,omitempty"` // tricking the json marshalling engine
	Authors     Authors `json:"authors"`
	Volumes     Volumes `json:"volumes,omitempty"`
}

type SerieScoped struct {
	Serie
	Volumes VolumesScoped `json:"volumes,omitempty"`
}

type SerieGet struct {
	Serie
	Volumes []*Book `json:"volumes,omitempty"`
}

type SerieGetScoped struct {
	Serie
	Volumes []*BookScoped `json:"volumes,omitempty"`
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
	gob.Register(Flash{})
}
