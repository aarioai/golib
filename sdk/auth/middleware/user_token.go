package middleware

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/httpsvr/response"
	"github.com/aarioai/golib/sdk/auth"
	"github.com/aarioai/golib/sdk/irisz"
	"github.com/aarioai/golib/typez"
	"github.com/kataras/iris/v12"
)

func (w *Middleware) parseUserAuth(ictx iris.Context) (svc typez.Svc, uid, vuid uint64, e *ae.Error) {
	var ok bool
	if svc, uid, vuid, ok = irisz.Uid(ictx); ok {
		return
	}
	svc, uid, vuid, _, _, e = w.parser(ictx)
	if e != nil {
		return
	}
	if ok = irisz.SetUid(ictx, svc, uid, vuid); !ok {
		e = NewE("middleware set uid failed")
	}
	return
}

// TryLoadUserAuth
// 这里是middleware handler不能有返回值
func (w *Middleware) TryLoadUserAuth(ictx iris.Context) {
	defer ictx.Next()
	if w.configSection == "" && w.parser == nil {
		Panic("requires config section or parser")
	}
	if w.parser == nil {
		auth.New(w.app, w.configSection).LoadUserAuth(ictx)
		return
	}
	w.parseUserAuth(ictx)
}

// MustLoadUserAuth
// 这里是middleware handler不能有返回值
func (w *Middleware) MustLoadUserAuth(ictx iris.Context) {
	if w.configSection == "" && w.parser == nil {
		Panic("middleware requires config section or parser")
	}
	var e *ae.Error
	if w.parser == nil {
		_, _, _, e = auth.New(w.app, w.configSection).LoadUserAuth(ictx)
	} else {
		_, _, _, e = w.parseUserAuth(ictx)
	}
	if e != nil {
		response.JsonE(ictx, e)
		return
	}
	ictx.Next()
}

func (w *Middleware) parseUserAuthorization(ictx iris.Context) (svc typez.Svc, uid, vuid uint64, ttl int64, accessToken string, e *ae.Error) {
	var ok bool
	if svc, uid, vuid, ttl, accessToken, ok = irisz.UserAuthorization(ictx); ok {
		return
	}
	svc, uid, vuid, ttl, accessToken, e = w.parser(ictx)
	if e != nil {
		return
	}
	if ok = irisz.SetUserAuthorization(ictx, svc, uid, vuid, ttl, accessToken); !ok {
		e = NewE("middleware set uid failed")
	}
	return
}

// TryLoadUserAuthorization try to load user authorization and set them into iris context values
// it relies on each request, and the token contains each client user agent. every user agents' token are all different
// header Authorization -> query access_token -> header access_token/AccessToken/X-AccessToken -> cookie access_token
// 这里是middleware handler不能有返回值
func (w *Middleware) TryLoadUserAuthorization(ictx iris.Context) {
	defer ictx.Next()
	if w.configSection == "" && w.parser == nil {
		Panic("requires config section or parser")
	}
	if w.parser == nil {
		auth.New(w.app, w.configSection).LoadUserAuthorization(ictx)
		return
	}
	w.parseUserAuthorization(ictx)
}

func (w *Middleware) MustLoadUserAuthorization(ictx iris.Context) {
	if w.configSection == "" && w.parser == nil {
		Panic("requires config section or parser")
	}
	var e *ae.Error
	if w.parser == nil {
		_, _, _, _, _, e = auth.New(w.app, w.configSection).LoadUserAuthorization(ictx)
	} else {
		_, _, _, _, _, e = w.parseUserAuthorization(ictx)
	}
	if e != nil {
		response.JsonE(ictx, e)
		return
	}
	ictx.Next()
}
