package openid

import "github.com/aarioai/airis-driver/driver/index"

type AccessKey struct {
	Id uint `db:"id" json:"id"`
}

func (t AccessKey) Table() string {
	return ""
}
func (t AccessKey) Indexes() index.Indexes {
	return index.NewIndexes(
		index.Primary("id"),
	)
}
