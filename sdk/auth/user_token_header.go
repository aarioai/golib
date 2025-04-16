package auth

import (
	"github.com/aarioai/airis/aa/httpsvr/request"
	"github.com/aarioai/golib/enumz"
	"strings"
)

func ApiAccessToken(r *request.Request) string {
	// Authorization: Bearer $access_token
	s := r.HeaderFast(enumz.HeaderAuthorization)
	if s == "" {
		s = r.QueryWildFast(enumz.ParamAccessToken)
	}
	if s == "" {
		return ""
	}
	p := strings.IndexByte(s, ' ')
	if p == 0 {
		return s
	}
	return s[p+1:]
}

func AccessToken(r *request.Request) string {
	if cookie, err := r.Cookie(enumz.ParamAccessToken); err == nil {
		return cookie.Value
	}
	return ApiAccessToken(r)
}

// SseAccessToken 由于原生sse不允许传递header
func SseAccessToken(r *request.Request) string {
	return AccessToken(r)
}

func ViewAccessToken(r *request.Request) string {
	if cookie, err := r.Cookie(enumz.ParamAccessToken); err == nil {
		return cookie.Value
	}
	return ""
}

func ViewOnLogout(r *request.Request) bool {
	if cookie, err := r.Cookie(enumz.ParamLogout); err == nil {
		return cookie.Value == "1"
	}
	return false
}
