package entity

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis/aa/atype"
)

type Object struct {
	Id        uint           `db:"id"`
	V         string         `db:"v"`
	Name      string         `db:"name"`
	CreatedAt atype.Datetime `db:"created_at"`
	UpdatedAt atype.Datetime `db:"updated_at"`
}

func (t Object) Table() string {
	return "casbin_object"
}

func (t Object) Indexes() index.Indexes {
	return index.NewIndexes(
		index.Primary("id"),
		index.Unique("v"),
	)
}
