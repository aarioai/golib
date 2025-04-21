package auth

import (
	"fmt"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/sdk/auth/cache"
	"github.com/aarioai/golib/sdk/auth/configz"
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
		CheckConfig()
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

func panicOnEmpty(name, s string) {
	if s != "" {
		return
	}
	panic(fmt.Sprintf("libsdk_auth: configz.%s not set", name))
}
func CheckConfig() {
	panicOnEmpty("UserTokenCryptMd5Key", configz.UserTokenCryptMd5Key)
	panicOnEmpty("UserTokenShuffleBase", configz.UserTokenShuffleBase)
	coding.ValidateShuffleEncryptKeys(configz.UserTokenShuffleBase)
}
