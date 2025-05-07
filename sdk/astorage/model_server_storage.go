package astorage

import (
	"context"
	"fmt"
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/golib/sdk/astorage/entity"
)

type cos struct {
	g *AStorage
}

func newCos(g *AStorage) *cos {
	return &cos{g: g}
}

func (m *cos) add(ctx context.Context, t entity.ServiceStorage) (string, *ae.Error) {
	db := m.g.db()

	qs := fmt.Sprintf(`
		INSERT INTO %s
			SET k=?, v=?, remark=?, created_at=now, updated_at=now()
		ON DUPLICATE KEY UPDATE v=?, remark=?, status=0, updated_at=now()
`, t.Table())
	e := db.Exec(ctx, qs, t.K, t.V, t.Remark, t.V, t.Remark)
	return t.K, e
}

func (m *cos) updateStatus(ctx context.Context, k string, status aenum.Status) *ae.Error {
	db := m.g.db()

	t := entity.ServiceStorage{
		K:      k,
		Status: status,
	}
	qs := fmt.Sprintf(`
		UPDATE %s SET status=?, updated_at=now() WHERE k=? LIMIT 1
`, t.Table())
	e := db.Exec(ctx, qs, t.Status, t.K)
	return e
}
func (m *cos) del(ctx context.Context, k string) *ae.Error {
	return m.updateStatus(ctx, k, aenum.Deleted)
}

func (m *cos) find(ctx context.Context, k string) (entity.ServiceStorage, *ae.Error) {
	db := m.g.db()

	t := entity.ServiceStorage{K: k}
	ps := fmt.Sprintf(`
		SELECT v, remark, created_at, updated_at FROM %s WHERE k=? AND status>=0 LIMIT 1
	`, t.Table())
	e := db.ScanX(ctx, ps, k, &t.V, &t.Remark, &t.CreatedAt, &t.UpdatedAt)
	return t, e
}

func (m *cos) all(ctx context.Context) ([]entity.ServiceStorage, *ae.Error) {
	db := m.g.db()

	t := entity.ServiceStorage{}
	qs := fmt.Sprintf(`
		SELECT v, remark, created_at, updated_at FROM %s
	`, t.Table())

	rows, e := db.Query(ctx, qs)
	if e != nil {
		return nil, e
	}

	defer rows.Close()
	var err error
	cfgs := make([]entity.ServiceStorage, 0)
	for rows.Next() {
		t = entity.ServiceStorage{}
		if err = rows.Scan(&t.V, &t.Remark, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, driver.NewMysqlError(err)
		}
		cfgs = append(cfgs, t)
	}
	if len(cfgs) == 0 {
		return nil, ae.ErrorNoRowsAvailable
	}
	return cfgs, nil
}
