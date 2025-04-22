package apppay

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/golib/sdk/weixin/weixinpay"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
	"sync"
	"time"
)

const (
	AppType = weixinpay.AtrApp
)

type Service struct {
	app *aa.App
	loc *time.Location

	Appid      string
	PayService *weixinpay.Service
}

var (
	services sync.Map
)

// New
// 一个APP，可以关联任意个支付mch账号
func New(app *aa.App, appid string, config weixinpay.Config) (*Service, error) {
	sk := appid + config.Mchid
	var s *Service
	sv, ok := services.Load(sk)
	if ok {
		if s, ok = sv.(*Service); ok && s != nil {
			return s, nil
		}
		services.Delete(sk)
	}
	if appid == "" {
		return nil, fmt.Errorf("appid is empty")
	}
	payservice, err := weixinpay.New(app, config)
	if err != nil {
		return nil, err
	}
	s = &Service{
		app:        app,
		loc:        app.Config.TimeLocation,
		Appid:      appid,
		PayService: payservice,
	}
	services.LoadOrStore(sk, s)
	return s, nil
}

func (s *Service) NewError(msg string, a ...any) error {
	msg = afmt.Sprintf(msg, a...)
	return fmt.Errorf("weixin %s pay (appid:%s, mchid:%s): %s", AppType, s.Appid, s.PayService.Config.Mchid, msg)
}

func (s *Service) client(ctx context.Context) (*app.AppApiService, error) {
	client, err := s.PayService.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &app.AppApiService{Client: client}, nil
}
