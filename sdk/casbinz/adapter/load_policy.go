package adapter

import (
	"context"
	"fmt"
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/airis/pkg/arrmap"
	"github.com/aarioai/golib/sdk/casbinz/adapter/entity"
	"github.com/aarioai/golib/sdk/casbinz/enum"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"strconv"
)

var (
	actMap = map[string]string{
		"r":          enum.Read.String(),
		"w":          enum.Write.String(),
		"rw":         enum.ReadWrite.String(),
		"r|w":        enum.ReadWrite.String(),
		"wr":         enum.ReadWrite.String(),
		"w|r":        enum.ReadWrite.String(),
		"read|write": enum.ReadWrite.String(),
	}
)

func handleAct(act string) string {
	if v, ok := actMap[act]; ok {
		return v
	}
	return act
}

// LoadPolicy loads all policy rules from the storage.
func (a *Adapter) LoadPolicy(model model.Model) error {
	ctx := context.Background()
	roles, e := a.selectAllRoles(ctx)
	if e != nil {
		return handleDriverError(e)
	}
	if err := a.loadPolicies(ctx, model, roles); err != nil {
		return err
	}
	return a.loadGroups(ctx, model, roles)
}

// p, role, object, act
func (a *Adapter) loadPolicies(ctx context.Context, m model.Model, roles map[uint]entity.Role) error {
	objs, e := a.selectAllObjects(ctx)
	if e != nil {
		return handleDriverError(e)
	}
	policies, e := a.selectAllPolicies(ctx)
	if e != nil {
		return handleDriverError(e)
	}
	for _, policy := range policies {
		policyObjects := policy.Objects.Uints()
		if len(policyObjects) == 0 {
			continue
		}

		ptype := afmt.DefaultIfZero(policy.Ptype, DefaultPolicyPtype)
		role, ok := roles[policy.Role]
		if !ok {
			return NewAdapterError("policy:%d role:%d not found in roles table", policy.Id, policy.Role)
		}

		for _, policyObject := range policyObjects {
			obj, ok := objs[policyObject]
			if !ok {
				return NewAdapterError("policy:%d object:%d not found in object table", policy.Id, policyObject)
			}
			act := handleAct(policy.Act)
			rule := []string{ptype, role.V0, obj.V, act}
			if err := persist.LoadPolicyArray(rule, m); err != nil {
				return err
			}
		}

	}
	return nil
}

// inherit parent role
func inheritUserRoles(userId uint64, userRoles []uint, roles map[uint]entity.Role) ([]uint, error) {
	extRoles := make([]uint, 0)
	for _, userRole := range userRoles {
		r, ok := roles[userRole]
		if !ok {
			return nil, NewAdapterError("user:%d role:%d not found in roles table", userId, userRole)
		}
		extRoles = append(extRoles, userRole)
		if r.Pid == 0 {
			continue
		}
		if _, ok = roles[r.Pid]; !ok {
			return nil, NewAdapterError("user:%d role:%d pid:%d not found in roles table", userId, userRole, r.Pid)
		}
		extRoles = append(extRoles, r.Pid)
	}

	return arrmap.Compact(extRoles, false), nil
}

// g, uid, role
func (a *Adapter) loadGroups(ctx context.Context, m model.Model, roles map[uint]entity.Role) error {
	users, e := a.selectAllUsers(ctx)
	if e != nil {
		return handleDriverError(e)
	}
	var err error
	for _, user := range users {
		userRoles := user.Roles.Uints()
		if len(userRoles) == 0 {
			continue
		}

		userRoles, err = inheritUserRoles(user.Id, userRoles, roles)
		if err != nil {
			return err
		}

		for _, userRole := range userRoles {
			r, _ := roles[userRole]
			ptype := afmt.DefaultIfZero(r.Ptype, DefaultGroupPtype)
			userSid := strconv.FormatUint(user.Id, 10)
			rule := []string{ptype, userSid, r.V0}
			if err = persist.LoadPolicyArray(rule, m); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *Adapter) selectAllRoles(ctx context.Context) (map[uint]entity.Role, *ae.Error) {
	db := a.db()
	t := entity.Role{}
	qs := fmt.Sprintf(`
		SELECT id, pid, v0, ptype, effective_at, expire_at FROM %s
	`, t.Table())
	rows, e := db.Query(ctx, qs)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	roles := make(map[uint]entity.Role)
	for rows.Next() {
		var r entity.Role
		err := rows.Scan(&r.Id, &r.Pid, &r.V0, &r.Ptype, &r.EffectiveAt, &r.ExpireAt)
		if err != nil {
			return nil, driver.NewMysqlError(err)
		}
		roles[r.Id] = r
	}
	if len(roles) == 0 {
		return nil, ae.ErrorNoRowsAvailable
	}
	return roles, nil
}

func (a *Adapter) selectAllObjects(ctx context.Context) (map[uint]entity.Object, *ae.Error) {
	db := a.db()
	t := entity.Object{}
	qs := fmt.Sprintf(`
		SELECT id, v FROM %s
	`, t.Table())
	rows, e := db.Query(ctx, qs)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	objects := make(map[uint]entity.Object)
	for rows.Next() {
		var r entity.Object
		err := rows.Scan(&r.Id, &r.V)
		if err != nil {
			return nil, driver.NewMysqlError(err)
		}
		objects[r.Id] = r
	}
	if len(objects) == 0 {
		return nil, ae.ErrorNoRowsAvailable
	}
	return objects, nil
}

func (a *Adapter) selectAllPolicies(ctx context.Context) (map[uint]entity.Policy, *ae.Error) {
	db := a.db()
	t := entity.Policy{}
	qs := fmt.Sprintf(`
		SELECT id, ptype, role, objects, act FROM %s
	`, t.Table())
	rows, e := db.Query(ctx, qs)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	objects := make(map[uint]entity.Policy)
	for rows.Next() {
		var r entity.Policy
		err := rows.Scan(&r.Id, &r.Ptype, &r.Role, &r.Objects, &r.Act)
		if err != nil {
			return nil, driver.NewMysqlError(err)
		}
		objects[r.Id] = r
	}
	if len(objects) == 0 {
		return nil, ae.ErrorNoRowsAvailable
	}
	return objects, nil
}

func (a *Adapter) selectAllUsers(ctx context.Context) ([]entity.User, *ae.Error) {
	db := a.db()
	t := entity.User{}
	qs := fmt.Sprintf(`
		SELECT uid, roles FROM %s
	`, t.Table())
	rows, e := db.Query(ctx, qs)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	objects := make([]entity.User, 0)
	for rows.Next() {
		var r entity.User
		err := rows.Scan(&r.Id, &r.Roles)
		if err != nil {
			return nil, driver.NewMysqlError(err)
		}
		objects = append(objects, r)
	}
	if len(objects) == 0 {
		return nil, ae.ErrorNoRowsAvailable
	}
	return objects, nil
}
