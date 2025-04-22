package weixinpay

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"net/http"
)

func (s *Service) getRsaNotifyHandler() (*notify.Handler, error) {
	mtx.RLock()
	if s.rsaNotifyHandler != nil {
		mtx.RUnlock()
		return s.rsaNotifyHandler, nil
	}
	mtx.RUnlock()
	mtx.Lock()
	defer mtx.Unlock()

	// 2. 获取商户号对应的微信支付平台证书访问器
	certVisitor := downloader.MgrInstance().GetCertificateVisitor(s.Config.Mchid)
	// 3. 使用证书访问器初始化 `notify.Handler`
	handler, err := notify.NewRSANotifyHandler(s.Config.MchApiV3Key, verifiers.NewSHA256WithRSAVerifier(certVisitor))
	if err != nil {
		return nil, s.NewError("notify rsa error: %s", err.Error())
	}
	s.rsaNotifyHandler = handler
	return handler, nil
}

// https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_5.shtml
func (s *Service) HandlePayNotify(ctx context.Context, r *http.Request) (Transaction, error) {
	trans := &payments.Transaction{}
	handler, err := s.getRsaNotifyHandler()
	if err != nil {
		return Transaction{}, err
	}
	// https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_5.shtml
	notifyResp, err := handler.ParseNotifyRequest(ctx, r, trans) // 如果验签未通过，或者解密失败
	if err != nil {
		return Transaction{}, s.NewError("parse pay notify request error: %s", err.Error())
	}
	if notifyResp.EventType != "TRANSACTION.SUCCESS" {
		return Transaction{}, s.NewError(notifyResp.EventType + " " + notifyResp.Summary)
	}

	t, err := ToTransaction(*trans, notifyResp.ID)
	return t, err
}
