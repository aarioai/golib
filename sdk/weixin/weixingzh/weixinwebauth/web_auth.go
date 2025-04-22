package weixinwebauth

import (
	"errors"
	"github.com/aarioai/airis/pkg/httpc"
	"github.com/aarioai/golib/sdk/weixin/weixingzh/base"
	"html/template"
	"strings"
)

const (
	GrantAccessTokenUrl = "https://api.weixin.qq.com/sns/oauth2/access_token"
	GetUserInfoUrl      = "https://api.weixin.qq.com/sns/userinfo"
)

// 用户网页授权
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
// 1 第一步：用户同意授权，获取code
// 2 第二步：通过code换取网页授权access_token
// 3 第三步：刷新access_token（如果需要）
// 4 第四步：拉取用户信息(需scope为 snsapi_userinfo)

// Auth 第二步：通过code换取网页授权access_token
// 这个是特殊的网页授权access token (grant type: authorization_code)，不同于 weixingzh 里的基础API access token (grant type: client_credential)
func (s *Service) Auth(code string) (AuthResult, error) {
	params := map[string]string{
		"appid":      s.appid,
		"secret":     s.secret,
		"code":       code,
		"grant_type": "authorization_code",
	}
	_, body, err := httpc.Get(GrantAccessTokenUrl, params)
	if err != nil {
		return AuthResult{}, err
	}
	var result AuthResult
	_, err = base.ParseResult(body, &result)
	return result, err
}

func parseHeadImage(url template.URL) template.URL {
	if url == "" {
		return ""
	}
	// 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），不修改的话，默认是132像素
	// https://thirdwx.qlogo.cn/mmopen/vi_32/ajNVdqHZLLA9a7MzbibMezia2OJRAYZOYeicAI8NaYPEh0mMrxsBySA0zq7ficumlc4kosbppPSL3iaWeh3soW0zeYA/132
	p := strings.LastIndexByte(string(url), '/')
	return url[0:p] + "/0"
}

func (s *Service) UserInfo(auth AuthResult) (UserInfo, error) {
	if auth.Scope != "snsapi_userinfo" {
		return UserInfo{}, errors.New("weixin web auth get userinfo, scope must be snsapi_userinfo")
	}

	params := map[string]string{
		"access_token": auth.AccessToken,
		"openid":       auth.Openid,
		"lang":         "zh_CN",
	}

	_, body, err := httpc.Get(GetUserInfoUrl, params)
	if err != nil {
		return UserInfo{}, err
	}
	var result UserInfoResult
	_, err = base.ParseResult(body, &result)

	result.HeadImg = parseHeadImage(result.HeadImg)
	return result.UserInfo, nil
}

func (s *Service) Code2UserInfo(code string) (UserInfo, error) {
	auth, err := s.Auth(code)
	if err != nil {
		return UserInfo{}, err
	}
	return s.UserInfo(auth)
}
