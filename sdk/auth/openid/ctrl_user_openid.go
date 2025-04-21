package openid

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/aarioai/golib/sdk/auth/dtoz"
	"github.com/aarioai/golib/sdk/auth/midiris"
	"github.com/kataras/iris/v12"
)

func (s *Service) GrantUserOpenid(ictx iris.Context) {
	defer ictx.Next()
	r, resp, ctx := httpsvr.New(ictx)
	defer resp.CloseWith(r)

	svc, uid, _, ok := midiris.Uid(ictx, false)
	if !ok {
		resp.WriteE(ae.ErrorUnauthorized)
		return
	}

	openid, ttl, e := s.EncodeUserOpenid(ctx, svc, uid)
	if e != nil {
		resp.WriteE(e)
		return
	}

	data := dtoz.UserOpenidResponse{
		Openid:    openid,
		ExpiresIn: atype.ToDurationSeconds(ttl),
	}
	resp.Write(data)
}
