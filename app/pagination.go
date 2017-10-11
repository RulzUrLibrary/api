package app

import (
	"github.com/paul-bismuth/library/utils"
)

const DEFAULT_LIMIT = 10

type Meta struct {
	Limit  int `query:"limit" json:"limit" validate:"gt=0,lte=50"`
	Offset int `query:"offset" json:"offset" validate:"gt=0"`
	Count  int `json:"count,omitempty"`
}

func NewMeta() Meta {
	return Meta{Limit: DEFAULT_LIMIT, Count: -1}
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
