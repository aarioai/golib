package casbinz

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/sdk/casbinz/adapter"
)

func (a *adapter.Adapter) Init(ctx acontext.Context) *ae.Error {
	a.CreateTables()
}

func (a *adapter.Adapter) CreateTables() *ae.Error {

}
