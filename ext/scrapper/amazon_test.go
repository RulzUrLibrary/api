package scrapper

import (
	"fmt"
	"testing"
)

type Result struct {
	title  string
	serie  string
	number int
}

func (r Result) String() string {
	return fmt.Sprintf(`{"%s" "%s" "%d"}`, r.title, r.serie, r.number)
}

func (r Result) Equal(_r Result) bool {
	return r.title == _r.title && r.serie == _r.serie && r.number == _r.number
}

type testPairs struct {
	title  string
	result Result
}

var tests = []testPairs{
	{
		"The Lean Startup: How Constant Innovation Creates Radically Successful Businesses",
		Result{"The Lean Startup: How Constant Innovation Creates Radically Successful Businesses", "", 0},
	},
	{"FullMetal Alchemist Vol.1", Result{"", "FullMetal Alchemist", 1}},
	{"FullMetal Alchemist - tome 08 (8)", Result{"", "FullMetal Alchemist", 8}},
	{"Fullmetal Alchemist, Tome 12", Result{"", "Fullmetal Alchemist", 12}},
	{"Lanfeust Odyssey T08 - Tseu-Hi la gardienne", Result{"Tseu-Hi la gardienne", "Lanfeust Odyssey", 8}},
	{"Dune – Tome 1", Result{"", "Dune", 1}},
	{"Dune (1)", Result{"", "Dune", 1}},
}

func TestGetTitle(t *testing.T) {
	var result Result
	for _, pair := range tests {
		result.title, result.serie, result.number = getTitle(pair.title)
		if result != pair.result {
			t.Error(
				"For", pair.title,
				"expected", pair.result,
				"got", result,
			)
		}
	}
}
