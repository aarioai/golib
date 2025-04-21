package mmc

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/sdk/auth/cache"
	"sync"
	"time"
)

type Service struct {
	app                 *aa.App
	loc                 *time.Location
	h                   *cache.Cache
	pubDERBase64KeyName string // rsa public key in config key name
	privDERKeyName      string // rsa private key in config key name
	gcmKeyName          string // gcm secret in config key name
}

var (
	once sync.Once
	s    *Service
)

func New(app *aa.App, redisConfigSection, pubDERBase64KeyName, privDERConfigKeyName, gcmKeyName string) *Service {
	once.Do(func() {
		s = &Service{app: app,
			loc:                 app.Config.TimeLocation,
			h:                   cache.New(app, redisConfigSection),
			pubDERBase64KeyName: pubDERBase64KeyName,
			privDERKeyName:      privDERConfigKeyName,
			gcmKeyName:          gcmKeyName,
		}
	})
	return s
}

func NewE(format string, args ...any) *ae.Error {
	return ae.NewE("libsdk_auth_mmc: "+format, args...)
}

func NewError(err error) *ae.Error {
	if err == nil {
		return nil
	}
	return ae.NewE("libsdk_auth_mmc: " + err.Error())
}
