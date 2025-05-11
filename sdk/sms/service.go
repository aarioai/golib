package sms

import (
	"context"
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/sdk/sms/aliyun"
	"github.com/aarioai/golib/sdk/sms/cache"
	"sync"
	"time"
)

const prefix = "libsdk_sms: "

type Service struct {
	app        *aa.App
	loc        *time.Location
	aliyun     *aliyun.Aliyun
	h          *cache.Cache
	mongo      *mongodb.Model
	enableLog  bool
	initSignal chan struct{}
}

var (
	once sync.Once
	s    *Service
)

func New(app *aa.App, redisConfigSection string) *Service {
	once.Do(func() {
		ca := cache.New(app, redisConfigSection)
		s = &Service{app: app,
			loc:        app.Config.TimeLocation,
			h:          ca,
			initSignal: make(chan struct{}, 1),
		}
		go s.checkInitReady()
	})
	return s
}

func (s *Service) checkInitReady() {
	ticker := time.NewTimer(time.Second * 5)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.app.Log.Warn(context.Background(), prefix+"not init yet")
		case <-s.initSignal:
			return
		}
	}
}

func (s *Service) WithMongo(mongo *mongodb.Model) *Service {
	s.mongo = mongo
	s.enableLog = true
	return s
}

func (s *Service) WithAliyun(accessKey, accessSecret string, regionId ...string) *Service {
	s.aliyun = aliyun.NewAliyun(accessKey, accessSecret, regionId...)
	return s
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
