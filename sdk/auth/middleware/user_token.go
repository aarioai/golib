package middleware

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/httpsvr/response"
	"github.com/aarioai/golib/sdk/auth"
	"github.com/aarioai/golib/sdk/auth/midiris"
	"github.com/aarioai/golib/typez"
	"github.com/kataras/iris/v12"
)

func (w *Middleware) parseUserAuth(ictx iris.Context) (svc typez.Svc, uid, vuid uint64, e *ae.Error) {
	var ok bool
	if svc, uid, vuid, ok = midiris.Uid(ictx); ok {
		return
	}
	svc, uid, vuid, _, _, e = w.parser(ictx)
	if e != nil {
		return
	}
	if ok = midiris.SetUid(ictx, svc, uid, vuid); !ok {
		e = NewE("middleware set uid failed")
	}
	return
}

func (w *Middleware) TryLoadUserAuth(ictx iris.Context) (svc typez.Svc, uid, vuid uint64, e *ae.Error) {
	defer ictx.Next()
	if w.configSection == "" && w.parser == nil {
		panic(prefix + "middleware requires config section or parser")
	}
	if w.parser == nil {
		return auth.New(w.app, w.configSection).LoadUserAuth(ictx)
	}
	return w.parseUserAuth(ictx)
}

func (w *Middleware) MustLoadUserAuth(ictx iris.Context) (svc typez.Svc, uid, vuid uint64) {
	var e *ae.Error
	if svc, uid, vuid, e = w.parseUserAuth(ictx); e != nil {
		response.JsonE(ictx, e)
		return
	}
	ictx.Next()
	return
}

func (w *Middleware) parseUserAuthorization(ictx iris.Context) (svc typez.Svc, uid, vuid uint64, ttl int64, accessToken string, e *ae.Error) {
	var ok bool
	if svc, uid, vuid, ttl, accessToken, ok = midiris.UserAuthorization(ictx); ok {
		return
	}
	svc, uid, vuid, ttl, accessToken, e = w.parser(ictx)
	if e != nil {
		return
	}
	if ok = midiris.SetUserAuthorization(ictx, svc, uid, vuid, ttl, accessToken); !ok {
		e = NewE("middleware set uid failed")
	}
	return
}

// TryLoadUserAuthorization try to load user authorization and set them into iris context values
// it relies on each request, and the token contains each client user agent. every user agents' token are all different
// header Authorization -> query access_token -> header access_token/AccessToken/X-AccessToken -> cookie access_token
func (w *Middleware) TryLoadUserAuthorization(ictx iris.Context) (svc typez.Svc, uid, vuid uint64, ttl int64, accessToken string, e *ae.Error) {
	defer ictx.Next()
	if w.configSection == "" && w.parser == nil {
		panic(prefix + "middleware requires config section or parser")
	}
	if w.parser == nil {
		return auth.New(w.app, w.configSection).LoadUserAuthorization(ictx)
	}
	return w.parseUserAuthorization(ictx)
}

func (w *Middleware) MustLoadUserAuthorization(ictx iris.Context) (svc typez.Svc, uid, vuid uint64, ttl int64, accessToken string) {
	var e *ae.Error
	if svc, uid, vuid, ttl, accessToken, e = w.TryLoadUserAuthorization(ictx); e != nil {
		response.JsonE(ictx, e)
		return
	}
	ictx.Next()
	return
}
