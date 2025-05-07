package entity

import (
	"encoding/json"
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/aa/atype"
)

// ClientStorage 客户端保存的配置文件
// 简单就是美，用json字段，既简单、又可以纠错
type ClientStorage struct {
	Uid       uint64          `db:"uid"`
	K         string          `db:"k"`
	V         json.RawMessage `db:"v"`
	Readonly  uint8           `db:"readonly"`
	Remark    string          `db:"remark"`
	Status    aenum.Status    `db:"status"`
	CreatedAt atype.Datetime  `db:"created_at"`
	UpdatedAt atype.Datetime  `db:"updated_at"`
}

func (t ClientStorage) Table() string {
	return "astorage_client"
}

func (t ClientStorage) Indexes() index.Indexes {
	return index.NewIndexes(
		index.Primary("uid", "k"),
	)
}
