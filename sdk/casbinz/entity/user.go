package entity

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis/aa/atype"
)

type User struct {
	Id        uint64          `db:"id"`
	Roles     atype.NullUints `db:"roles"`
	CreatedAt atype.Datetime  `db:"created_at"`
	UpdatedAt atype.Datetime  `db:"updated_at"`
}

func (t User) Table() string {
	return "casbin_user"
}

func (t User) Indexes() index.Indexes {
	return index.NewIndexes(
		index.Primary("id"),
		index.Index("role"),
	)
}
