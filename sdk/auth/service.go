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

const prefix = "libsdk_auth: "

type Service struct {
	app *aa.App
	loc *time.Location
	h   *cache.Cache
}

var (
	once sync.Once
	s    *Service
)

func New(app *aa.App, redisConfigSection string) *Service {
	once.Do(func() {
		s = &Service{app: app,
			loc: app.Config.TimeLocation,
			h:   cache.New(app, redisConfigSection),
		}
	})
	return s
}

func NewCode(code int, format string, args ...any) *ae.Error {
	return ae.New(code, afmt.Sprintf(prefix+format, args...))
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

func panicOnEmpty(name, s string) {
	if s != "" {
		return
	}
	panic(fmt.Sprintf(prefix+"configz.%s not set", name))
}
