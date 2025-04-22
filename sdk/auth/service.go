package auth

import (
	"fmt"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/golib/sdk/auth/cache"
	"sync"
	"time"
)

type Service struct {
	app      *aa.App
	loc      *time.Location
	h        *cache.Cache
	withVuid bool
}

var (
	once sync.Once
	s    *Service
)

func New(app *aa.App, redisConfigSection string, withVuid bool) *Service {
	once.Do(func() {
		s = &Service{app: app,
			loc:      app.Config.TimeLocation,
			h:        cache.New(app, redisConfigSection),
			withVuid: withVuid,
		}
	})
	return s
}

func NewCode(code int, format string, args ...any) *ae.Error {
	return ae.New(code, afmt.Sprintf("libsdk_auth: "+format, args...))
}

func NewE(format string, args ...any) *ae.Error {
	return ae.NewE("libsdk_auth: "+format, args...)
}

func NewError(err error) *ae.Error {
	if err == nil {
		return nil
	}
	return ae.NewE("libsdk_auth: " + err.Error())
}

func panicE(msg string, e *ae.Error) {
	panic("libsdk_auth: " + msg + " " + e.Text())
}
func panicMsg(msg string, args ...any) {
	panic(afmt.Sprintf("libsdk_auth: "+msg, args...))
}

func panicOnEmpty(name, s string) {
	if s != "" {
		return
	}
	panic(fmt.Sprintf("libsdk_auth: configz.%s not set", name))
}
