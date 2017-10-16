package app

import (
	"github.com/rulzurlibrary/api/utils"
	"strconv"
)

const DEFAULT_LIMIT = 10

type Meta struct {
	Limit  int `query:"limit" json:"limit" validate:"gt=0,lte=50"`
	Offset int `query:"offset" json:"offset" validate:"gte=0"`
	Count  int `json:"count,omitempty"`
}

type Pagination struct {
	Page  int `query:"page" validate:"gt=0"`
	Count int `query:"-"`
}

func NewMeta() Meta {
	return Meta{DEFAULT_LIMIT, 0, -1}
}

func NewPagination() Pagination {
	return Pagination{1, -1}
}

func (p Pagination) Offset() int {
	return (p.Page - 1) * p.Limit()
}

func (p Pagination) Limit() int {
	return DEFAULT_LIMIT
}

func (p Pagination) Last() bool {
	return p.Count != -1 && p.Page == utils.Ceil(p.Count, p.Limit())
}

func (p Pagination) Current() string {
	return strconv.Itoa(p.Page)
}

func (p Pagination) Pages() []string {
	page := p.Page
	pages := []string{
		"...", strconv.Itoa(page - 2), strconv.Itoa(page - 1),
		strconv.Itoa(page), strconv.Itoa(page + 1),
		strconv.Itoa(page + 2), "...",
	}
	if page < 5 {
		pages = pages[5-page:]
	}
	pages = append([]string{"1"}, pages...)

	if p.Count != -1 {
		last := utils.Ceil(p.Count, p.Limit())
		if page > last-5 {
			pages = pages[:len(pages)-4+last-page]
		}
		pages = append(pages, strconv.Itoa(last))
	}
	return pages
}
