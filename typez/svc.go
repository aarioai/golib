package typez

import (
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
)

type Svc atype.Uint24

// Valid 最大值：36^4-1 = 1679615  --> 采用base36编码，4位占位符 1679615
func (s Svc) Valid() bool {
	return s > 0 && s < 1679615
}

func (s Svc) String() string {
	return types.FormatUint(uint64(s))
}

func (s Svc) Uint32() uint32 { return uint32(s) }

func (s Svc) Or(defaultSvc Svc) Svc {
	if s.Valid() {
		return s
	}
	return defaultSvc
}
