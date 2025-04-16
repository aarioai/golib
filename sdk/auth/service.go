package auth

import (
	"fmt"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/sdk/auth/cache"
	"github.com/aarioai/golib/sdk/auth/configz"
	"sync"
	"time"
)

type Service struct {
	app *aa.App
	loc *time.Location
	h   *cache.Cache
}

var (
	once     sync.Once
	instance *Service
)

func New(app *aa.App, redisConfigSection string) *Service {
	once.Do(func() {
		CheckConfig()
		ca := cache.New(app, redisConfigSection)
		instance = &Service{app: app,
			loc: app.Config.TimeLocation,
			h:   ca,
		}
	})
	return instance
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
}
