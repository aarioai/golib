package weixinwebauth

import (
	"errors"
	"html/template"
)

type Response struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}

func (r Response) CheckError() error {
	if r.AccessToken == "" {
		return errors.New("parse auth result failed")
	}
	return nil
}

type UserInfo struct {
	Unionid   string       `json:"unionid"`
	Openid    string       `json:"openid"`
	Nickname  string       `json:"nickname"`
	HeadImg   template.URL `json:"headimgurl"`
	Privilege []string     `json:"privilege"`
}

// UserInfoResult
// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
type UserInfoResult struct {
	UserInfo

	// 2021年10月起，微信已不再提供 sex, country, province, city, language 信息，都返回默认值
	// https://developers.weixin.qq.com/community/develop/doc/00028edbe3c58081e7cc834705b801?blockType=1
	Sex      any `json:"sex"` //
	Language any `json:"language"`
	Country  any `json:"country"`
	Province any `json:"province"`
	City     any `json:"city"`
}

func (r UserInfoResult) CheckError() error {
	if r.Openid == "" {
		return errors.New("parse userinfo result failed")
	}
	return nil
}
