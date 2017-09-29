package db

import (
	"database/sql"
	"fmt"
	"github.com/paul-bismuth/library/utils"
)

type Serie struct {
	id          int
	name        sql.NullString
	title       sql.NullString
	serie       sql.NullString
	description sql.NullString
	volumes     Volumes
	authors     Authors
}

func (s Serie) ToStructs() *utils.Serie {
	isbn := ""
	volumes := s.volumes.ToStructs()

	if len(volumes) > 0 && volumes[0].Number == 0 {
		isbn = volumes[0].Isbn
		volumes = nil
	}
	return &utils.Serie{
		s.id, s.name.String, s.description.String, s.title.String, isbn,
		s.authors.ToStructs(), volumes,
	}
}

type Volume struct {
	isbn   string
	number sql.NullInt64
	owned  bool
}
type Volumes []Volume

func (v Volume) String() string {
	return fmt.Sprintf("%d", v.number.Int64)
}

func (v Volumes) ToStructs() (volumes utils.Volumes) {
	var prev string
	var volume Volume

	for _, volume = range v {
		if volume.isbn != prev {
			volumes = append(
				volumes,
				utils.Volume{int(volume.number.Int64), volume.isbn, volume.owned},
			)
		}
		prev = volume.isbn
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

func (s *Series) Gets() (series []*utils.Serie) {
	for _, id := range s.order {
		series = append(series, s.series[id].ToStructs())
	}
	return series
}
