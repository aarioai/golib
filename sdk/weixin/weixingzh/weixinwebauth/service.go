package weixinwebauth

import (
	"github.com/aarioai/airis/aa"
)

type Service struct {
	app    *aa.App
	appid  string
	secret string
}

func New(app *aa.App, appid string, secret string) *Service {
	return &Service{
		app:    app,
		appid:  appid,
		secret: secret,
	}
}
