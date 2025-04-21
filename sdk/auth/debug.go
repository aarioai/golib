package auth

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/aa/httpsvr/request"
	"github.com/aarioai/airis/aa/httpsvr/response"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/sdk/auth/configz"
	"github.com/kataras/iris/v12"
	"time"
)

func parseCookies(ictx iris.Context) map[string]any {
	cookies := map[string]any{
		"server_time":   atype.Now(s.loc),
		"time_location": s.app.Config.TimeLocation.String(),
	}
	ctx := acontext.FromIris(ictx)
	ictx.VisitAllCookies(func(k, v string) {
		cookies[k] = v
		if enumz.ParamAccessToken != k {
			return
		}
		svc, uid, vuid, ua, psid, authAt, expiresIn, factor, _, e := s.DbgDecryptUserToken(ctx, v)
		if e != nil {
			cookies["decrypt_token_error"] = e.Text()
		} else {
			cookies["decrypt_token"] = map[string]any{
				"svc":        svc,
				"uid":        uid,
				"vuid":       vuid,
				"UA":         ua.Name(),
				"psid":       psid,
				"auth_at":    time.Unix(authAt, 0).In(s.app.Config.TimeLocation).Format(s.app.Config.TimeFormat),
				"expires_in": time.Unix(authAt, 0).Add(time.Duration(expiresIn) * time.Second).In(s.app.Config.TimeLocation).Format(s.app.Config.TimeFormat),
				"factor":     factor,
			}
		}
	})
	return cookies
}

// DebugMyCookies 用于微信H5、小程序等调试
func (s *Service) DebugMyCookies(ictx iris.Context) {
	defer ictx.Next()
	debugToken := request.QueryWild(ictx, enumz.ParamDebugToken)
	if debugToken != configz.DebugToken {
		response.JsonE(ictx, ae.ErrorPageExpired)
		return
	}
	cookies := parseCookies(ictx)
	response.JSON(ictx, cookies)
}

func (s *Service) DebugMyCookiesJSONP(ictx iris.Context) {
	defer ictx.Next()
	callback := request.QueryWild(ictx, enumz.ParamCallback)
	debugToken := request.QueryWild(ictx, enumz.ParamDebugToken)
	if debugToken != configz.DebugToken {
		response.JsonE(ictx, ae.ErrorPageExpired)
		return
	}
	cookies := parseCookies(ictx)
	if callback == "" {
		callback = enumz.ParamCallback
	}
	ictx.JSONP(cookies, iris.JSONP{Callback: callback})
}
