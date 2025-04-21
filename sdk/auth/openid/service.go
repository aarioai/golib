package openid

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/golib/typez"
	"sync"
	"time"
)

type Service struct {
	app *aa.App
	loc *time.Location
	//h             *cache.Cache
	secretHandler func(ctx context.Context, app *aa.App, svc typez.Svc) (appid, secret string, e *ae.Error)
}

var (
	once sync.Once
	s    *Service
)

func New(app *aa.App, secretHandler func(ctx context.Context, app *aa.App, svc typez.Svc) (string, string, *ae.Error)) *Service {
	once.Do(func() {
		s = &Service{app: app,
			loc: app.Config.TimeLocation,
			//h:             cache.New(app, redisConfigSection),
			secretHandler: secretHandler,
		}
	})
	return s
}

func NewCode(code int, format string, args ...any) *ae.Error {
	return ae.New(code, afmt.Sprintf("libsdk_auth_openid: "+format, args...))
}

func NewE(format string, args ...any) *ae.Error {
	return ae.NewE("libsdk_auth_openid: "+format, args...)
}

func NewError(err error) *ae.Error {
	if err == nil {
		return nil
	}
	return ae.NewE("libsdk_auth_openid: " + err.Error())
}

func panicOnEmpty(name, s string) {
	if s != "" {
		return
	}
	panic(fmt.Sprintf("libsdk_auth_openid: configz.%s not set", name))
}
