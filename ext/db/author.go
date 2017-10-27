package db

import (
	"bytes"
	"fmt"
	"github.com/lib/pq"
	"github.com/rulzurlibrary/api/utils"
	"strconv"
)

type Author utils.Author

func (a *Author) Scan(src interface{}) (err error) {
	var elems [][]byte

	if elems, err = parseRow(src.([]byte), []byte{','}); err != nil {
		return
	}

	if len(elems) != 2 {
		return fmt.Errorf("element is not a valid author")
	}
	if a.Id, err = strconv.ParseUint(string(elems[0]), 10, 64); err != nil {
		return
	}

	a.Name = string(elems[1])
	return
}

type Authors struct {
	*utils.Authors
}

func (a *Authors) Scan(src interface{}) error {
	authors := []Author{}
	if bytes.Equal(src.([]byte), []byte(`{"(,)"}`)) {
		a.Authors = &utils.Authors{}
		return nil
	}
	if err := pq.Array(&authors).Scan(src); err != nil {
		return err
	}
	a.Authors = &utils.Authors{}
	for _, author := range authors {
		*a.Authors = append(*a.Authors, utils.Author(author))
	}
	return nil
}
