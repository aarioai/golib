package adapter

import (
	"context"
	"fmt"
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/sdk/casbinz/entity"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"strconv"
	"strings"
)

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
		role, ok := roles[policy.Role]
		if !ok {
			return NewAdapterError("policy:%d role:%d not found in roles table", policy.Id, policy.Role)
		}

		var object strings.Builder
		for i, policyObject := range policyObjects {
			obj, ok := objs[policyObject]
			if !ok {
				return NewAdapterError("policy:%d object:%d not found in object table", policy.Id, policyObject)
			}
			if i > 0 {
				object.WriteString(",")
			}
			object.WriteString(obj.V)
		}

		ptype := policy.Ptype
		if ptype == "" {
			ptype = "p"
		}
		rule := []string{ptype, role.V0, object.String(), policy.Act}
		if err := persist.LoadPolicyArray(rule, m); err != nil {
			return err
		}
	}
	return nil
}

// g, uid, role
func (a *Adapter) loadGroups(ctx context.Context, m model.Model, roles map[uint]entity.Role) error {
	users, e := a.selectAllUsers(ctx)
	if e != nil {
		return handleDriverError(e)
	}
	for _, user := range users {
		userRoles := user.Roles.Uints()
		if len(userRoles) == 0 {
			continue
		}
		for _, userRole := range userRoles {
			r, ok := roles[userRole]
			if !ok {
				return NewAdapterError("user:%d role:%d not found in roles table", user.Id, userRole)
			}
			ptype := r.Ptype
			if ptype == "" {
				ptype = "g"
			}
			userSid := strconv.FormatUint(user.Id, 10)
			rule := []string{ptype, userSid, r.V0}
			if err := persist.LoadPolicyArray(rule, m); err != nil {
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
