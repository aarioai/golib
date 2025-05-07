package astorage

import (
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/aarioai/golib/sdk/astorage/entity"
	"github.com/kataras/iris/v12"
)

func (g *AStorage) getConfigOfService(ictx iris.Context) {
	defer ictx.Next()
	r, resp, ctx := httpsvr.New(ictx)
	k, e0 := r.QueryString("k")
	if e := resp.FirstError(e0); e != nil {
		return
	}

	m := newCos(g)
	cfg, e := m.find(ctx, k)
	if e != nil {
		resp.WriteE(e)
		return
	}

	resp.Write(toCosDto(cfg))
}

func (g *AStorage) getConfigsOfService(ictx iris.Context) {
	defer ictx.Next()
	_, resp, ctx := httpsvr.New(ictx)
	m := newCos(g)
	cfgs, e := m.all(ctx)
	if e != nil {
		resp.WriteE(e)
		return
	}
	resp.Write(toCosDtoes(cfgs))
}
func (g *AStorage) deleteConfigOfService(ictx iris.Context) {
	defer ictx.Next()
	r, resp, ctx := httpsvr.New(ictx)
	k, e0 := r.QueryString("k")
	if e := resp.FirstError(e0); e != nil {
		return
	}
	m := newCos(g)
	e := m.del(ctx, k)
	resp.WriteE(e)
}
func (g *AStorage) postConfigOfService(ictx iris.Context) {
	defer ictx.Next()
	r, resp, ctx := httpsvr.New(ictx)
	k, e0 := r.BodyString("k")
	v, e1 := r.BodyBytes("v")
	remark, e2 := r.Body("remark", false)
	if e := resp.FirstError(e0, e1, e2); e != nil {
		return
	}

	m := newCos(g)
	t := entity.ServiceStorage{
		K:      k,
		V:      v,
		Remark: remark.String(),
		Status: 0,
	}
	if _, e := m.add(ctx, t); e != nil {
		resp.WriteE(e)
		return
	}
	resp.WriteAliasId("k", t.K)
}
