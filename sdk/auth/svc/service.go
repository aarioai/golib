package svc

import (
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis-driver/driver/mysqli"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/golib/sdk/auth/cache"
	"sync"
	"time"
)

type Service struct {
	app                *aa.App
	loc                *time.Location
	h                  *cache.Cache
	mysqlConfigSection string
}

var (
	once sync.Once
	s    *Service
)

func New(app *aa.App, redisConfigSection, mysqlConfigSection string) *Service {
	once.Do(func() {
		s = &Service{app: app,
			loc:                app.Config.TimeLocation,
			h:                  cache.New(app, redisConfigSection),
			mysqlConfigSection: mysqlConfigSection,
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

func (s *Service) DB() *mysqli.DB {
	return mysqli.NewDriver(driver.NewMysqlPool(s.app, s.mysqlConfigSection))
}
