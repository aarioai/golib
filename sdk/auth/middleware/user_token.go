package middleware

import (
	"github.com/aarioai/airis/aa/httpsvr/response"
	"github.com/aarioai/golib/sdk/auth"
	"github.com/kataras/iris/v12"
)

// TryLoadUserAuthorization try to load user authorization and set them into iris context values
// it relies on each request, and the token contains each client user agent. every user agents' token are all different
// header Authorization -> query access_token -> header access_token/AccessToken/X-AccessToken -> cookie access_token
func (w *Middleware) TryLoadUserAuthorization(ictx iris.Context) {
	defer ictx.Next()
	auth.New(w.app, w.authRedisSection).LoadUserAuthorization(ictx)
}

func (w *Middleware) MustLoadUserAuthorization(ictx iris.Context) {
	_, _, _, e := auth.New(w.app, w.authRedisSection).LoadUserAuthorization(ictx)
	if e != nil {
		response.JsonE(ictx, e)
		return
	}
	ictx.Next()
}
