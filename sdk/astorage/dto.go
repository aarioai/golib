package astorage

import (
	"encoding/json"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/golib/sdk/astorage/entity"
)

type simpleCocDto struct {
	K string          `json:"k"`
	V json.RawMessage `json:"v"`
}
type cocDto struct {
	K         string          `json:"k"`
	Uid       uint64          `json:"uid"`
	V         json.RawMessage `json:"v"`
	Readonly  uint8           `json:"readonly"`
	Remark    string          `json:"remark"`
	CreatedAt atype.Datetime  `json:"created_at"`
	UpdatedAt atype.Datetime  `json:"updated_at"`
}

type cosDto struct {
	K         string          `json:"k"`
	V         json.RawMessage `json:"v"`
	Remark    string          `json:"remark"`
	Status    aenum.Status    `json:"status"`
	CreatedAt atype.Datetime  `json:"created_at"`
	UpdatedAt atype.Datetime  `json:"updated_at"`
}

func toSimpleCocDtoes(cs []entity.ClientStorage) []simpleCocDto {
	a := make([]simpleCocDto, len(cs))
	for i, c := range cs {
		a[i].K = c.K
		a[i].V = c.V
	}
	return a
}
func toCocDto(c entity.ClientStorage) cocDto {
	return cocDto{
		K:         c.K,
		Uid:       c.Uid,
		V:         c.V,
		Readonly:  c.Readonly,
		Remark:    c.Remark,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCosDto(c entity.ServiceStorage) cosDto {
	return cosDto{
		K:         c.K,
		V:         c.V,
		Remark:    c.Remark,
		Status:    c.Status,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCosDtoes(cs []entity.ServiceStorage) []cosDto {
	a := make([]cosDto, len(cs))
	for i, c := range cs {
		a[i] = toCosDto(c)
	}
	return a
}
