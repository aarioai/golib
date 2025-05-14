package adapter

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/sdk/casbinz/adapter/entity"
)

var (
	DefaultPolicyPtype = "p"
	DefaultGroupPtype  = "g"
)

func (a *Adapter) Init(ctx context.Context) {
	ae.PanicOn(a.createTables(ctx))
}

func (a *Adapter) createTables(ctx context.Context) *ae.Error {
	db := a.db()
	var tableN int
	obj := entity.Object{}.Table()
	policy := entity.Policy{}.Table()
	role := entity.Role{}.Table()
	user := entity.User{}.Table()
	expectedTables := []any{obj, policy, role, user}
	qs := fmt.Sprintf(`
		SELECT COUNT(*) FROM information_schema.TABLES
			WHERE TABLE_SCHEMA='%s' AND TABLE_NAME IN (?, ?, ?, ?)
	`, db.Schema)
	e := db.ScanArgs(ctx, qs, expectedTables, &tableN)
	if e != nil {
		return e
	}

	if tableN == len(expectedTables) {
		return nil
	}
	a.app.Log.Info(ctx, prefix+"schema %s create tables: %s, %s, %s, %s", db.Schema, obj, policy, role, user)
	
	tx, e := db.Begin(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Recover()

	qsObj := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
		    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
		    v VARCHAR(255) NOT NULL,
		    name VARCHAR(255) NOT NULL DEFAULT '',
		    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  		updated_at DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
		    PRIMARY KEY (id),
	  		UNIQUE INDEX u_v (v)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`, obj)
	if e = tx.Exec(ctx, qsObj); e != nil {
		a.app.Check(ctx, tx.Rollback())
		return e
	}

	qsPolicy := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			ptype VARCHAR(6) NOT NULL DEFAULT '',
			role INT UNSIGNED NOT NULL DEFAULT 0,
			objects JSON DEFAULT NULL,
			act VARCHAR(15) NOT NULL DEFAULT '',
		    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  		updated_at DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
		    PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`, policy)
	if e = tx.Exec(ctx, qsPolicy); e != nil {
		a.app.Check(ctx, tx.Rollback())
		return e
	}

	qsRole := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			pid INT UNSIGNED NOT NULL DEFAULT 0,
			v0 VARCHAR(15) NOT NULL,
		    name VARCHAR(255) NOT NULL DEFAULT '',
			ptype VARCHAR(6) NOT NULL DEFAULT '',
  			effective_at DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
  			expire_at DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
		    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  		updated_at DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
		    PRIMARY KEY (id),
	  		UNIQUE INDEX u_v0 (v0)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`, role)
	if e = tx.Exec(ctx, qsRole); e != nil {
		a.app.Check(ctx, tx.Rollback())
		return e
	}

	qsUser := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INT UNSIGNED NOT NULL AUTO_INCREMENT,
			roles JSON DEFAULT NULL,
		    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  		updated_at DATETIME NOT NULL DEFAULT '0000-00-00 00:00:00',
		    PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`, user)
	if e = tx.Exec(ctx, qsUser); e != nil {
		a.app.Check(ctx, tx.Rollback())
		return e
	}
	return tx.Commit()
}
