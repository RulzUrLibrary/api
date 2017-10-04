package db

import (
	"database/sql"
	"github.com/paul-bismuth/library/utils"
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

func (s Serie) base() *utils.Serie {
	return &utils.Serie{
		s.id, s.name.String, s.description.String, s.title.String, "", "", nil,
		s.authors.ToStructs(), nil,
	}
}

func (s Serie) ToSerie() (serie *utils.Serie) {
	serie = s.base()
	serie.Volumes = s.volumes.ToVolumes()
	if len(serie.Volumes) > 0 && serie.Volumes[0].Number == 0 {
		book := serie.Volumes[0]
		serie.Isbn = book.Isbn
		serie.Title = book.Title
		serie.Volumes = nil
	}
	return
}

func (s Serie) ToSerieScoped() (serie *utils.SerieScoped) {
	serie = &utils.SerieScoped{Serie: *s.base()}
	serie.Volumes = s.volumes.ToVolumesScoped()
	if len(serie.Volumes) > 0 && serie.Volumes[0].Number == 0 {
		book := serie.Volumes[0]
		serie.Isbn = book.Isbn
		serie.Title = book.Title
		serie.Volumes = nil
	}
	return
}

func (s Serie) ToSerieGet() (serie *utils.SerieGet) {
	serie = &utils.SerieGet{Serie: *s.base()}
	serie.Volumes = s.volumes.ToBooks()
	if len(serie.Volumes) > 0 && serie.Volumes[0].Number == 0 {
		book := serie.Volumes[0]
		serie.Isbn = book.Isbn
		serie.Title = book.Title
		serie.Thumbnail = book.Thumbnail
		serie.Description = book.Description
		serie.Volumes = nil
	}
	return
}

func (s Serie) ToSerieGetScoped() (serie *utils.SerieGetScoped) {
	serie = &utils.SerieGetScoped{Serie: *s.base()}
	serie.Volumes = s.volumes.ToBooksScoped()
	if len(serie.Volumes) > 0 && serie.Volumes[0].Number == 0 {
		book := serie.Volumes[0]
		serie.Isbn = book.Isbn
		serie.Title = book.Title
		serie.Thumbnail = book.Thumbnail
		serie.Description = book.Description
		serie.Owned = &book.Owned
		serie.Volumes = nil
	}
	return
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

func (s *Series) ToSeries() (series []*utils.Serie) {
	for _, id := range s.order {
		series = append(series, s.series[id].ToSerie())
	}
	return
}

func (s *Series) ToSeriesScoped() (series []*utils.SerieScoped) {
	for _, id := range s.order {
		series = append(series, s.series[id].ToSerieScoped())
	}
	return
}
func (s *Series) ToSeriesGet() (series []*utils.SerieGet) {
	for _, id := range s.order {
		series = append(series, s.series[id].ToSerieGet())
	}
	return
}

func (s *Series) ToSeriesGetScoped() (series []*utils.SerieGetScoped) {
	for _, id := range s.order {
		series = append(series, s.series[id].ToSerieGetScoped())
	}
	return
}
