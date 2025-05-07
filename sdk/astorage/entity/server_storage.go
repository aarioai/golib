package entity

import (
	"encoding/json"
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/aa/atype"
)

// 第三方服务端配置
// 简单就是美，用json字段，既简单、又可以纠错
type ServiceStorage struct {
	K         string          `db:"k"`
	V         json.RawMessage `db:"v"`
	Remark    string          `db:"remark"`
	Status    aenum.Status    `db:"status"`
	CreatedAt atype.Datetime  `db:"created_at"`
	UpdatedAt atype.Datetime  `db:"updated_at"`
}

func (t ServiceStorage) Table() string {
	return "astorage_service"
}

func (t ServiceStorage) Indexes() index.Indexes {
	return index.NewIndexes(
		index.Primary("k"),
	)
}
