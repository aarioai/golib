package midiris

import (
	"github.com/aarioai/airis/aa/httpsvr/request"
	"github.com/aarioai/golib/enumz"
	"github.com/kataras/iris/v12"
	"strings"
)

// AccessToken
// 中间件层，尽量不要使用 request.Request（可能会导致后期进入controller层，multipart等被清空）
// header Authorization -> query access_token -> header access_token/AccessToken/X-AccessToken -> cookie access_token
// 由于原生SSE和git hook等不允许传递header，因此支持多种传值方式会更加通用
func AccessToken(ictx iris.Context) string {
	s := ictx.GetHeader(enumz.HeaderAuthorization)
	if s == "" {
		s = request.QueryWild(ictx, enumz.ParamAccessToken) // include header and cookie
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

func ViewOnLogout(r *request.Request) bool {
	if cookie, err := r.Cookie(enumz.ParamLogout); err == nil {
		return cookie.Value == "1"
	}
	return false
}
