package db

import (
	"database/sql"
	"github.com/ixday/echo-hello/utils"
)

type Author struct {
	id   sql.NullInt64
	name sql.NullString
}
type Authors []Author

func (a Authors) ToStructs() utils.Authors {
	authors := utils.Authors{}
	ids := map[int64]bool{}

	for _, author := range a {
		if _, ok := ids[author.id.Int64]; !ok && author.id.Valid {
			ids[author.id.Int64] = true
			authors = append(authors, &utils.Author{Name: author.name.String})
		}
	}
	return authors
}
