package typez

import (
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
)

type Svc atype.Uint24

func (s Svc) Valid() bool {
	return s > 0
}

func (s Svc) String() string {
	return types.FormatUint(uint64(s))
}

func (s Svc) Or(defaultSvc Svc) Svc {
	if s.Valid() {
		return s
	}
	return defaultSvc
}
