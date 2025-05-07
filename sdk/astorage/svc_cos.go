package astorage

import (
	"context"
	"encoding/json"
	"github.com/aarioai/airis/aa/ae"
)

func (g *AStorage) GetConfigValueOfService(ctx context.Context, k string) (json.RawMessage, *ae.Error) {
	m := newCos(g)
	cfg, e := m.find(ctx, k)
	if e != nil {
		return nil, e
	}
	return cfg.V, nil
}
