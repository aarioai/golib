package svc

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/golib/sdk/auth/configz"
)

type Svc struct {
	Id        Id             `db:"id" json:"id"`
	Sid       Sid            `db:"sid" json:"sid"`
	Name      string         `db:"name" json:"name"`
	Logo      atype.Image    `db:"logo" json:"logo"`
	Iconfont  string         `db:"iconfont" json:"iconfont"`
	Status    aenum.Status   `db:"status" json:"status"`
	CreatedAt atype.Datetime `db:"created_at" json:"created_at"`
	UpdatedAt atype.Datetime `db:"updated_at" json:"updated_at"`
}

func (t Svc) Table() string {
	return configz.DbPrefix + "_svc"
}
func (t Svc) Indexes() index.Indexes {
	return index.NewIndexes(
		index.Primary("id"),
		index.Unique("sid"),
		index.FullText("name"),
	)
}
