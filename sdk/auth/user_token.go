package auth

import (
	"context"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/sdk/auth/configz"
	"github.com/aarioai/golib/sdk/auth/dtoz"
	"github.com/aarioai/golib/sdk/auth/midiris"
	"github.com/aarioai/golib/typez"
	"github.com/kataras/iris/v12"
	"time"
)

type CredentialsType uint8

const (
	DefaultCredentials CredentialsType = 0
	ApiCredentials     CredentialsType = 1
	ViewCredentials    CredentialsType = 2
	SseCredentials     CredentialsType = 3
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

// ParseUserAuthorization
func (s *Service) ParseUserAuthorization(ictx iris.Context) (svc typez.Svc, uid, vuid uint64, authAt int64, atoken string, e *ae.Error) {
	if atoken = midiris.AccessToken(ictx); atoken == "" {
		e = ae.ErrorUnauthorized
		return
	}
	ctx := ictx.Request().Context()
	di := midiris.DeviceInfo(ictx)
	var psid string
	var factor int64
	var ua enumz.UA
	svc, uid, vuid, ua, psid, authAt, factor, _, e = s.decryptUserToken(ctx, atoken)
	if e != nil {
		return
	}

	if !ua.Is(di.UA) || psid != di.PSID {
		s.app.Log.Notice(ctx, "user token %s UA(%s) %s != UA(%s) %s or psid:%s != %s", atoken, ua.String(), ua.Name(), di.UA.String(), di.UA.Name(), psid, di.PSID)
		e = ae.ErrorUnauthorized
		return
	}

	cachedFactor, ok := s.h.LoadUserTokenFactor(ctx, svc, uid, ua)
	if !ok || factor != cachedFactor {
		e = ae.ErrorLoginTimeout
		return
	}

	now := time.Now().Unix()
	expiresAt := authAt + configz.UserTokenTTLs
	if now > expiresAt {
		e = ae.ErrorLoginTimeout
		return
	}
	return
}

// LoadUserAuthorization parse user authorization then set them into context
func (s *Service) LoadUserAuthorization(ictx iris.Context) (svc typez.Svc, uid, vuid uint64, e *ae.Error) {
	var ok bool
	if svc, uid, vuid, ok = midiris.Uid(ictx, s.withVuid); ok {
		return
	}
	svc, uid, vuid, _, _, e = s.ParseUserAuthorization(ictx)
	if e != nil {
		return
	}
	if ok = midiris.SetUid(ictx, svc, uid, vuid); !ok {
		return 0, 0, 0, NewE("iris middleware set uid failed")
	}
	return
}

func (s *Service) UserLogout(ctx context.Context, svc typez.Svc, uid uint64, ua enumz.UA) {
	// 废除之前的 factor
	s.h.IncrUserTokenFactor(ctx, svc, uid, ua)
}
