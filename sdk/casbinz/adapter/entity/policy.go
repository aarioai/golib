package entity

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis/aa/atype"
)

type Policy struct {
	Id        uint            `db:"id"`
	Ptype     string          `db:"ptype"`
	Role      uint            `db:"role"`
	Objects   atype.NullUints `db:"objects"`
	Act       string          `db:"name"`
	CreatedAt atype.Datetime  `db:"created_at"`
	UpdatedAt atype.Datetime  `db:"updated_at"`
}

func (t Policy) Table() string {
	return "casbin_policy"
}

func (t Policy) Indexes() index.Indexes {
	return index.NewIndexes(
		index.Primary("id"),
	)
}
