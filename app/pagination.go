package app

import (
	"github.com/paul-bismuth/library/utils"
	"strconv"
)

type Meta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count,omitempty"`
}

type Pagination struct {
	Limit  limit  `query:"limit"`
	Offset offset `query:"offset"`
}

type offset int
type limit int

func (o *offset) UnmarshalParam(param string) error {
	i, err := strconv.Atoi(param)
	if err != nil {
		return err
	}
	if i < 0 {
		return utils.ErrOffset
	}
	*o = offset(i)
	return nil
}

func (l *limit) UnmarshalParam(param string) error {
	i, err := strconv.Atoi(param)
	if err != nil {
		return err
	}
	if i < 0 || i > 50 {
		return utils.ErrLimit
	}

	*l = limit(i)
	return nil
}

func NewPagination() Pagination {
	return Pagination{Limit: 10}
}

func (m Meta) First() bool {
	return m.Offset == 0
}

func (m Meta) Last() bool {
	return m.Count <= m.Limit || m.Offset == (utils.Ceil(m.Count, m.Limit)-1)*m.Limit
}

func (m Meta) Current() int {
	if m.Limit == 0 {
		return 1
	}
	return (m.Offset / m.Limit) + 1
}

func (m Meta) Previous() int {
	return m.Current() - 1
}

func (m Meta) Next() int {
	return m.Current() + 1
}

func (m Meta) Pages() (pages []int) {
	for i := 0; i < utils.Ceil(m.Count, m.Limit); i++ {
		pages = append(pages, i+1)
	}
	return
}
