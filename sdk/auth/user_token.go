package auth

import (
	"context"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/httpsvr/request"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/sdk/auth/configz"
	"github.com/aarioai/golib/sdk/auth/dtoz"
	"github.com/aarioai/golib/typez"
	"time"
)

// NewUserToken
// secureLogin 是否是通过验证码等安全方式登录的
func (s *Service) NewUserToken(ctx context.Context, svc typez.Svc, uid, vuid uint64, ua enumz.UA, psid string, admin uint8, conflict, secureLogin bool) (*dtoz.Token, *ae.Error) {

	authAt := time.Now().Unix()
	factor, ok := s.h.IncrUserTokenFactor(ctx, svc, uid, ua)
	if !ok {
		return nil, NewE("incr user token factor failed")
	}

	atoken, e := s.encryptUserToken(svc, uid, vuid, ua, psid, authAt, factor, secureLogin)
	if e != nil {
		return nil, e
	}
	var rtoken string
	rtoken, e = s.encryptUserFreshToken(atoken)
	if e != nil {
		return nil, e
	}

	ui := configz.UserTokenTTLs
	t := packUserToken(secureLogin, atoken, rtoken, ui, configz.UserRefreshTokenTTLs, admin, nil, conflict)
	return &t, nil
}

func (s *Service) LoadUserCredential(ctx context.Context, r *request.Request, di typez.DeviceInfo) (atoken string, svc typez.Svc, uid, vuid uint64, authAt int64, e *ae.Error) {
	if atoken = ApiAccessToken(r); atoken == "" {
		e = ae.NewE("no credentials")
		return
	}

	var psid string
	var factor int64
	var ua enumz.UA
	svc, uid, vuid, ua, psid, authAt, factor, _, e = s.decryptUserToken(ctx, atoken)
	if e != nil {
		return
	}

	if !ua.Is(di.UA) || psid != di.PSID {
		e = ae.NewE("user token UA(%s) %s != UA(%s) %s or psid:%s != %s", ua.String(), ua.Name(), di.UA.String(), di.UA.Name(), psid, di.PSID)
	}

	cachedFactor, _ := s.h.LoadUserTokenFactor(ctx, svc, uid, ua)
	if factor != cachedFactor {
		e = ae.NewE("invalid access token factor")
	}
	return
}
