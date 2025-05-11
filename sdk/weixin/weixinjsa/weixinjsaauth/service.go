package weixinjsaauth

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
)

const prefix = "libsdk_weixinjsa: "

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

func NewE(format string, args ...any) *ae.Error {
	return ae.NewE(prefix+format, args...)
}

func NewError(err error) *ae.Error {
	if err == nil {
		return nil
	}
	return ae.NewE(prefix + err.Error())
}
