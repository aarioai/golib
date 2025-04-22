package weixintransfer

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/golib/sdk/weixin/weixinpay"
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"
	"sync"
	"time"
)

type Service struct {
	app *aa.App
	loc *time.Location

	Appid      string // 申请商户号的appid或商户号绑定的appid（企业号corpid即为此appid）
	PayService *weixinpay.Service
}

var (
	services sync.Map
)

// New
// 一个支付mch账号可以关联N个APP
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
	return fmt.Errorf("weixin transfer (appid:%s, mchid:%s): %s", s.Appid, s.PayService.Config.Mchid, msg)
}

func (s *Service) transferClient(ctx context.Context) (*transferbatch.TransferBatchApiService, error) {
	client, err := s.PayService.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &transferbatch.TransferBatchApiService{Client: client}, nil
}

func (s *Service) detailClient(ctx context.Context) (*transferbatch.TransferDetailApiService, error) {
	client, err := s.PayService.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &transferbatch.TransferDetailApiService{Client: client}, nil
}
