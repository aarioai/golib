package auth

import (
	"context"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/sdk/auth/configz"
	"github.com/aarioai/golib/sdk/auth/dtoz"
	"github.com/aarioai/golib/typez"
	"time"
)

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
