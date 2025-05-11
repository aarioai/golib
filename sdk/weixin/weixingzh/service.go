package weixingzh

import (
	"errors"
	"fmt"
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/golib/sdk/weixin/weixingzh/weixinwebauth"
	"github.com/redis/go-redis/v9"
)

const prefix = "libsdk_weixingzh: "

type Service struct {
	app    *aa.App
	appid  string
	secret string

	cacheCfgSection        string
	cacheClientCredential  string // cache key
	cacheJSSDKTicketPrefix string

	WebAuth *weixinwebauth.Service
}

func New(app *aa.App, appid string, secret string, redisCfgSection string) *Service {
	return &Service{
		app:    app,
		appid:  appid,
		secret: secret,

		cacheCfgSection:        redisCfgSection,
		cacheClientCredential:  fmt.Sprintf("libsdk:weixingzh:appid:%s:client_credential", appid),
		cacheJSSDKTicketPrefix: fmt.Sprintf("libsdk:weixingzh:appid:%s:js_ticket", appid),

		WebAuth: weixinwebauth.New(app, appid, secret),
	}
}

func (s *Service) rdb() (*redis.Client, error) {
	cli, e := driver.NewRedisPool(s.app, s.cacheCfgSection)
	if e != nil {
		return nil, errors.New(prefix + "load redis client failed: " + e.Text())
	}
	return cli, nil
}
