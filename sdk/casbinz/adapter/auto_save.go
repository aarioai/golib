package adapter

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/casbin/casbin/v2/model"
)

// 不支持自动保存，下面直接使用 errors.New("not implemented") 即可

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	return ae.ErrNotImplemented
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return ae.ErrNotImplemented
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return ae.ErrNotImplemented
}

// SavePolicy saves all policy rules to the storage.
func (a *Adapter) SavePolicy(model model.Model) error {
	return ae.ErrNotImplemented
}
