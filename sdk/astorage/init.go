package astorage

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/sdk/astorage/entity"
)

func (g *AStorage) Init() *ae.Error {
	return g.CreateTables()
}

func (g *AStorage) CreateTables() *ae.Error {
	_db := g.db()

	const expectedTableNum = 2
	ctx := context.Background()
	var (
		c entity.ClientStorage
		s entity.ServiceStorage

		tc int
	)
	// 每次启动都会检测，所以不要用到事务 FOR UPDATE
	ts := fmt.Sprintf(`
	SELECT COUNT(*) FROM information_schema.tables 
		WHERE TABLE_SCHEMA='%s' AND TABLE_NAME IN ('%s', '%s')
	`, _db.Schema, c.Table(), s.Table())
	if e := _db.ScanRow(ctx, ts, &tc); e != nil {
		return e
	}
	if expectedTableNum == tc {
		return nil
	}

	tx, e := _db.Begin(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Recover()

	cs := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		  uid BIGINT UNSIGNED NOT NULL DEFAULT 0,
		  k VARCHAR(255) NOT NULL DEFAULT '',
		  v JSON NOT NULL,
		  readonly TINYINT UNSIGNED NOT NULL DEFAULT 0,
		  remark VARCHAR(255) NOT NULL DEFAULT '',
		  status TINYINT NOT NULL DEFAULT 0,
		  created_at DATETIME NULL DEFAULT 0,
		  updated_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
		  PRIMARY KEY (uid, k)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`, c.Table())
	if e = tx.Exec(ctx, cs); e != nil {
		g.app.Check(ctx, tx.Rollback())
		return e
	}

	ss := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		  k VARCHAR(255) NOT NULL DEFAULT '',
		  v JSON NOT NULL,
		  remark VARCHAR(255) NOT NULL DEFAULT '',
		  status TINYINT NOT NULL DEFAULT 0,
		  created_at DATETIME NULL DEFAULT 0,
		  updated_at DATETIME NULL DEFAULT CURRENT_TIMESTAMP,
		  PRIMARY KEY (k)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`, s.Table())
	if e = tx.Exec(ctx, ss); e != nil {
		g.app.Check(ctx, tx.Rollback())
		return e
	}

	return tx.Commit()
}
