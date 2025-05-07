package astorage

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/aarioai/golib/sdk/astorage/entity"
	"github.com/kataras/iris/v12"
)

func (g *AStorage) getClientData(ictx iris.Context) {
	defer ictx.Next()
	r, resp, ctx := httpsvr.New(ictx)
	k, e0 := r.QueryString("k")
	uid, e1 := r.QueryUint64("uid")
	if e := resp.FirstError(e0, e1); e != nil {
		return
	}
	m := newCoc(g)
	cfg := entity.ClientStorage{K: k, Uid: uid}
	var e *ae.Error
	cfg, e = m.find(ctx, cfg.K, cfg.Uid)
	if e != nil {
		resp.WriteE(e)
		return
	}
	resp.Write(toCocDto(cfg))
}

func (g *AStorage) getClientDataByUid(ictx iris.Context) {
	defer ictx.Next()
	r, resp, ctx := httpsvr.New(ictx)
	uid, e0 := r.QueryUint64("uid")
	if e := resp.FirstError(e0); e != nil {
		return
	}
	m := newCoc(g)
	cfgs, e := m.findAllByUid(ctx, uid)
	if e != nil {
		resp.WriteE(e)
		return
	}
	resp.Write(toSimpleCocDtoes(cfgs))
}

func (g *AStorage) deleteClientData(ictx iris.Context) {
	defer ictx.Next()
	r, resp, ctx := httpsvr.New(ictx)
	k, e0 := r.Query("k")
	uid, e1 := r.QueryUint64("uid")
	if e := resp.FirstError(e0, e1); e != nil {
		return
	}
	m := newCoc(g)
	e := m.del(ctx, k.String(), uid)
	resp.WriteE(e)
}

func (g *AStorage) postClientData(ictx iris.Context) {
	defer ictx.Next()
	r, resp, ctx := httpsvr.New(ictx)
	k, e0 := r.Body("k")
	uid, e1 := r.QueryUint64("uid")
	v, e2 := r.Body("v")
	remark, e3 := r.Body("remark", false)
	if e := resp.FirstError(e0, e1, e2, e3); e != nil {
		return
	}
	m := newCoc(g)
	t := entity.ClientStorage{
		K:      k.String(),
		Uid:    uid,
		V:      v.Bytes(),
		Remark: remark.String(),
		Status: 0,
	}
	if _, _, e := m.add(ctx, t); e != nil {
		resp.WriteE(e)
		return
	}
	resp.WriteJointId("k", t.K, "uid", t.Uid)
}
