package sms

import (
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/golib/sdk/sms/aliyun"
	"github.com/aarioai/golib/sdk/sms/cache"
	"sync"
	"time"
)

type Service struct {
	app       *aa.App
	loc       *time.Location
	aliyun    *aliyun.Aliyun
	cache     *cache.Cache
	mongo     *mongodb.Model
	enableLog bool
}

var (
	instances sync.Map
)

func New(app *aa.App, redisConfigSection string) *Service {
	s, ok := instances.Load(redisConfigSection)
	if ok {
		if s != nil {
			return s.(*Service)
		}
		instances.Delete(redisConfigSection)
	}

	ca := cache.New(app, redisConfigSection)
	s = &Service{app: app,
		loc:   app.Config.TimeLocation,
		cache: ca,
	}
	s, _ = instances.LoadOrStore(redisConfigSection, s)
	ss := s.(*Service)
	ss.checkInitReady()
	return ss
}

func (s *Service) checkInitReady() {

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
