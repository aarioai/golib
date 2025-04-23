package svc

import (
	"context"
	"fmt"
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis/aa/ae"
)

func (s *Service) InsertSvc(ctx context.Context, t Svc) (Id, *ae.Error) {
	db := s.DB()
	qs := fmt.Sprintf(`
		INSERT INTO %s SET
			sid=?, name=?, logo=?, iconfont=?, status=?, created_at=now(), updated_at=now()
	`, t.Table())
	id, e := db.Insert(ctx, qs, t.Sid, t.Name, t.Logo, t.Iconfont, t.Status)
	return Id(id), e
}

func (s *Service) FindSvc(ctx context.Context, id Id) (Svc, *ae.Error) {
	db := s.DB()
	t := Svc{Id: id}
	qs := fmt.Sprintf(`SELECT sid, name, logo, iconfont, status, created_at, updated_at FROM %s WHERE id=? LIMIT 1`, t.Table())
	e := db.Scan(ctx, qs, uint64(id), &t.Sid, &t.Name, &t.Logo, &t.Iconfont, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	return t, e
}

func (s *Service) SeekSvc(ctx context.Context, sid Sid) (Svc, *ae.Error) {
	db := s.DB()
	t := Svc{Sid: sid}
	qs := fmt.Sprintf(`SELECT id, name, logo, iconfont, status, created_at, updated_at FROM %s WHERE sid=? LIMIT 1`, t.Table())
	e := db.ScanX(ctx, qs, string(sid), &t.Id, &t.Name, &t.Logo, &t.Iconfont, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	return t, e
}

func (s *Service) QuerySvc(ctx context.Context, name string) ([]Svc, *ae.Error) {
	db := s.DB()
	var t Svc
	qs := fmt.Sprintf(`SELECT id, sid, logo, iconfont, status, created_at, updated_at FROM %s WHERE name=?`, t.Table())
	rows, e := db.Query(ctx, qs, name)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	ts := make([]Svc, 0)
	for rows.Next() {
		n := Svc{Name: name}
		err := rows.Scan(&t.Id, &t.Sid, &t.Logo, &t.Iconfont, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, driver.NewMysqlError(err)
		}
		ts = append(ts, n)
	}
	if len(ts) == 0 {
		return nil, ae.ErrorNoRowsAvailable
	}
	return ts, e
}

func (s *Service) SearchSvc(ctx context.Context, name string) ([]Svc, *ae.Error) {
	db := s.DB()
	var t Svc
	qs := fmt.Sprintf(`SELECT id, sid, name, logo, iconfont, status, created_at, updated_at FROM %s WHERE name LIKE ?`, t.Table())
	rows, e := db.Query(ctx, qs, "%"+name+"%")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	ts := make([]Svc, 0)
	for rows.Next() {
		var n Svc
		err := rows.Scan(&t.Id, &t.Sid, &t.Name, &t.Logo, &t.Iconfont, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, driver.NewMysqlError(err)
		}
		ts = append(ts, n)
	}
	if len(ts) == 0 {
		return nil, ae.ErrorNoRowsAvailable
	}
	return ts, e
}
