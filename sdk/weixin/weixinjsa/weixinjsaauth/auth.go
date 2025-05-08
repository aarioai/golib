package weixinjsaauth

import (
	"context"
	"encoding/json"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/httpc"
)

const (
	GrantAccessTokenUrl = "https://api.weixin.qq.com/sns/jscode2session"
	GetUserInfoUrl      = "https://api.weixin.qq.com/sns/userinfo"
)

func (s *Service) Auth(ctx context.Context, code string) (Response, *ae.Error) {
	params := map[string]string{
		"appid":      s.appid,
		"secret":     s.secret,
		"js_code":    code,
		"grant_type": "authorization_code",
	}
	_, body, err := httpc.Get(ctx, GrantAccessTokenUrl, params)
	if err != nil {
		return Response{}, NewError(err)
	}
	var result ResponseWithError
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Response{}, NewError(err)
	}
	if result.Errcode != 0 {
		return Response{}, NewE("errcode: %d, errmsg: %s", result.Errcode, result.Errmsg)
	}
	return result.Response, nil
}
