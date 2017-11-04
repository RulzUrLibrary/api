package db

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/rulzurlibrary/api/ext/scrapper"
	"github.com/rulzurlibrary/api/utils"
	"strconv"
)

const sizeNotation = 3

type Notation utils.Notation

func (n *Notation) Scan(src interface{}) (err error) {
	var elems [][]byte
	var note float64

	if elems, err = parseRow(src.([]byte), []byte{','}); err != nil {
		return
	}

	if len(elems) != sizeNotation {
		return fmt.Errorf("element is not a valid notation")
	}
	if note, err = strconv.ParseFloat(string(elems[1]), 64); err != nil {
		return
	}
	n.Note = float32(note)
	n.Provider = string(elems[0])
	n.Link = string(elems[2])
	switch n.Provider {
	case scrapper.AMAZON_NAME:
		n.Max = scrapper.AMAZON_MAX_NOTATION
	case scrapper.SENSCRITIQUE_NAME:
		n.Max = scrapper.SENSCRITIQUE_MAX_NOTATION
	}
	return
}

type Notations struct {
	*utils.Notations
}

func (n *Notations) Scan(src interface{}) error {
	notations := []Notation{}
	n.Notations = &utils.Notations{}

	if emptyArray(sizeNotation, src) {
		return nil
	}
	if err := pq.Array(&notations).Scan(src); err != nil {
		return err
	}
	for _, notation := range notations {
		*n.Notations = append(*n.Notations, utils.Notation(notation))
	}
	return nil
}
