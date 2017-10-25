package db

import (
	"database/sql"
	"github.com/rulzurlibrary/api/utils"
)

type Serie struct {
	id          int
	name        sql.NullString
	title       sql.NullString
	serie       sql.NullString
	description sql.NullString
	volumes     *Books
	authors     Authors
}

func (s Serie) ToStructs(partial bool) *utils.Serie {
	serie := &utils.Serie{
		Id: s.id, Description: s.description.String,
		Authors: s.authors.ToStructs(), Volumes: s.volumes.ToStructs(partial),
	}
	if len(serie.Volumes) > 0 && serie.Volumes[0].Number == 0 {
		book := serie.Volumes[0]
		serie.Id = 0
		serie.Title = book.Title
		serie.Isbn = book.Isbn
		serie.Thumbnail = book.Thumbnail
		serie.Tags = book.Tags
		serie.Volumes = nil
	} else {
		serie.Name = s.name.String
	}
	return serie
}

type Series struct {
	series map[int]*Serie
	order  []int
}

func NewSeries() *Series {
	return &Series{make(map[int]*Serie), make([]int, 0)}
}

func (s *Series) Add(serie *Serie) {
	if _, ok := s.series[serie.id]; !ok {
		s.order = append(s.order, serie.id)
	}
	s.series[serie.id] = serie
}

func (s *Series) Get(id int) *Serie {
	return s.series[id]
}

func (s *Series) First() *Serie {
	if len(s.series) > 0 {
		return s.series[s.order[0]]
	}
	return nil
}

func (s *Series) ToStructs(partial bool) (series []*utils.Serie) {
	for _, id := range s.order {
		series = append(series, s.series[id].ToStructs(partial))
	}
	return
}
