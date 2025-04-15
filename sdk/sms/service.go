package sms

import (
	"github.com/aarioai/airis-driver/driver/mongodb"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/golib/sdk/sms/aliyun"
	"time"
)

type Service struct {
	app       *aa.App
	loc       *time.Location
	aliyun    *aliyun.Aliyun
	mongo     *mongodb.Model
	enableLog bool
}

func New(app *aa.App) *Service {
	s := &Service{app: app, loc: app.Config.TimeLocation}
	return s
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
