package astorage

import (
	"context"
	"fmt"
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/golib/sdk/astorage/entity"
)

type coc struct {
	g *AStorage
}

func newCoc(g *AStorage) *coc {
	return &coc{g: g}
}

// Checked and Set，强制以 Sid （不包括UID）来更新
func (m *coc) updateByKey(ctx context.Context, t entity.ClientStorage) *ae.Error {
	db := m.g.db()

	qs := fmt.Sprintf(`
		INSERT INTO %s
			SET uid=?, k=?, v=?, remark=?, created_at=now, updated_at=now()
		ON DUPLICATE KEY UPDATE v=?, remark=?, status=0, updated_at=now()
`, t.Table())
	e := db.Exec(ctx, qs, t.Uid, t.K, t.V, t.Remark, t.V, t.Remark)
	return e
}

func (m *coc) add(ctx context.Context, t entity.ClientStorage) (k string, uid uint64, e *ae.Error) {
	db := m.g.db()

	qs := fmt.Sprintf(`
		INSERT INTO %s
			SET uid=?, k=?, v=?, remark=?, created_at=now, updated_at=now()
		ON DUPLICATE KEY UPDATE v=?, remark=?, status=0, updated_at=now()
`, t.Table())
	e = db.Exec(ctx, qs, t.Uid, t.K, t.V, t.Remark, t.V, t.Remark)
	return t.K, t.Uid, e
}

func (m *coc) updateStatus(ctx context.Context, k string, uid uint64, status aenum.Status) *ae.Error {
	db := m.g.db()

	t := entity.ClientStorage{
		K:      k,
		Uid:    uid,
		Status: status,
	}
	qs := fmt.Sprintf(`
		UPDATE %s SET status=?, updated_at=now() WHERE k=? AND uid=? LIMIT 1
`, t.Table())
	e := db.Exec(ctx, qs, t.Status, t.K, t.Uid)
	return e
}

func (m *coc) del(ctx context.Context, k string, uid uint64) *ae.Error {
	return m.updateStatus(ctx, k, uid, aenum.Deleted)
}

func (m *coc) find(ctx context.Context, k string, uid uint64) (entity.ClientStorage, *ae.Error) {
	db := m.g.db()

	t := entity.ClientStorage{K: k}
	qs := fmt.Sprintf(`
		SELECT v, readonly, remark, created_at, updated_at FROM %s WHERE k=? AND uid=? AND status>=0 LIMIT 1
	`, t.Table())

	e := db.ScanArgs(ctx, qs, []any{k, uid}, &t.V, &t.Readonly, &t.Remark, &t.CreatedAt, &t.UpdatedAt)
	return t, e
}

func (m *coc) findAllByUid(ctx context.Context, uid uint64) ([]entity.ClientStorage, *ae.Error) {
	db := m.g.db()

	t := entity.ClientStorage{Uid: uid}
	ps := fmt.Sprintf(`
		SELECT k, v FROM %s WHERE uid=? AND status>=0
	`, t.Table())

	rows, e := db.Query(ctx, ps, uid)
	if e != nil {
		return nil, e
	}

	defer rows.Close()
	var err error
	cfgs := make([]entity.ClientStorage, 0)
	for rows.Next() {
		t = entity.ClientStorage{}
		if err = rows.Scan(&t.K, &t.V); err != nil {
			return nil, driver.NewMysqlError(err)
		}
		cfgs = append(cfgs, t)
	}
	if len(cfgs) == 0 {
		return nil, ae.ErrorNoRowsAvailable
	}
	// 一定要记得close rows，释放连接资源
	return cfgs, nil
}
