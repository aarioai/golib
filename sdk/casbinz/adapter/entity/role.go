package entity

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis/aa/atype"
)

type Role struct {
	Id          uint           `db:"id"`
	Pid         uint           `db:"pid"`
	V0          string         `db:"v0"`
	Name        string         `db:"name"`
	Ptype       string         `db:"ptype"`
	EffectiveAt atype.Datetime `db:"effective_at"`
	ExpireAt    atype.Datetime `db:"expire_at"`
	CreatedAt   atype.Datetime `db:"created_at"`
	UpdatedAt   atype.Datetime `db:"updated_at"`
}

func (t Role) Table() string {
	return "casbin_role"
}

func (t Role) Indexes() index.Indexes {
	return index.NewIndexes(
		index.Primary("id"),
		index.Unique("v0"),
	)
}
