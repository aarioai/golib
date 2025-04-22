package jsapay

import (
	"context"
	"github.com/aarioai/golib/sdk/weixin/weixinpay"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
	"time"
)

// Prepay 下单
// https://pay.weixin.qq.com/doc/v3/merchant/4012791857
func (s *Service) Prepay(ctx context.Context, t weixinpay.PrepayRequest) (PrepayDTO, error) {
	config := s.PayService.Config

	detail := parsePrepayDetail(t.Detail)
	sceneInfo := parsePrepaySceneInfo(t.SceneInfo)
	outTradeNo := base.NewOutTradeNo(t.OrderBatch, t.Total, time.Now().Format("150405"))

	pr := jsapi.PrepayRequest{
		Appid:         core.String(s.Appid),
		Mchid:         core.String(config.Mchid),
		Description:   core.String(t.Description),
		OutTradeNo:    core.String(outTradeNo),
		TimeExpire:    t.TimeExpire,
		Attach:        core.String(t.Attach),
		NotifyUrl:     core.String(config.NotifyUrl),
		GoodsTag:      core.String(t.GoodsTag),
		LimitPay:      t.LimitPay,
		SupportFapiao: core.Bool(t.Fapiao),
		Amount: &jsapi.Amount{
			Total:    base.FromMoney(t.Total),
			Currency: core.String(t.Currency.ISO4217()),
		},
		Payer:     &jsapi.Payer{Openid: core.String(t.PayerOpenid)},
		Detail:    detail,
		SceneInfo: sceneInfo,
		SettleInfo: &jsapi.SettleInfo{
			ProfitSharing: core.Bool(t.SettleProfitSharing),
		},
	}

	client, err := s.client(ctx)
	if err != nil {
		return PrepayDTO{}, err
	}

	resp, _, err := client.PrepayWithRequestPayment(ctx, pr)
	if err != nil {
		return PrepayDTO{}, s.NewError(err.Error())
	}
	prepayDTO := s.newPrepayDTO(*resp, outTradeNo)
	return prepayDTO, err
}

// CloseOrder 关闭订单
// https://pay.weixin.qq.com/doc/v3/merchant/4012791860
//
// 以下情况需要调用关单接口：
// 1. 商户订单支付失败需要生成新单号重新发起支付，要对原订单号调用关单，避免重复支付；
// 2. 系统下单后，用户支付超时，系统退出不再受理，避免用户继续，请调用关单接口。
func (s *Service) CloseOrder(ctx context.Context, outTradeNo string) error {
	config := s.PayService.Config
	client, err := s.client(ctx)
	if err != nil {
		return err
	}

	req := jsapi.CloseOrderRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(config.Mchid),
	}
	_, err = client.CloseOrder(ctx, req)
	return err
}

// QueryUnpaidOrder 未支付成功订单查询
func (s *Service) QueryUnpaidOrder(ctx context.Context, outTradeNo string) (weixinpay.Transaction, error) {
	config := s.PayService.Config
	client, err := s.client(ctx)
	if err != nil {
		return weixinpay.Transaction{}, err
	}

	req := jsapi.QueryOrderByOutTradeNoRequest{
		OutTradeNo: core.String(outTradeNo),
		Mchid:      core.String(config.Mchid),
	}
	resp, _, err := client.QueryOrderByOutTradeNo(ctx, req)
	if err != nil {
		return weixinpay.Transaction{}, err
	}

	return weixinpay.ToTransaction(*resp, "")
}

// QueryPaidOrder 支付成功订单查询
func (s *Service) QueryPaidOrder(ctx context.Context, transId string) (weixinpay.Transaction, error) {
	config := s.PayService.Config
	client, err := s.client(ctx)
	if err != nil {
		return weixinpay.Transaction{}, err
	}
	req := jsapi.QueryOrderByIdRequest{
		TransactionId: core.String(transId),
		Mchid:         core.String(config.Mchid),
	}
	resp, _, err := client.QueryOrderById(ctx, req)
	if err != nil {
		return weixinpay.Transaction{}, err
	}

	return weixinpay.ToTransaction(*resp, "")
}
