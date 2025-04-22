package weixinpay

import (
	"context"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

func (s *Service) handleRefundResult(resp *refunddomestic.Refund, result *core.APIResult, err error) (RefundResult,
	error) {
	if err != nil {
		return RefundResult{}, err
	}
	return toRefundResult(*resp, s.loc), nil
}

// Refund 退款申请
// https://pay.weixin.qq.com/doc/v3/merchant/4013070371
// 当交易发生之后一段时间内，由于买家或者卖家的原因需要退款时，卖家可以通过退款接口将支付款退还给买家，微信支付将在收到退款请求并且验证成功之后，按照退款规则将支付款按原路退到买家帐号上。
//
// 注意：
// 1、交易时间超过一年的订单无法提交退款
// 2、微信支付退款支持单笔交易分多次退款，多次退款需要提交原支付订单的商户订单号和设置不同的退款单号。申请退款总金额不能超过订单金额。 一笔退款失败后重新提交，请不要更换退款单号，请使用原商户退款单号
// 3、请求频率限制：150qps，即每秒钟正常的申请退款请求次数不超过150次，而调用失败报错时的频率限制为6QPS。
// 4、一笔订单最多支持50次部分退款（若需多次部分退款，请更换商户退款单号并间隔1分钟后再次调用）
// 5、如果同一个用户有多笔退款，建议分不同批次进行退款，避免并发退款导致退款失败
// 6、申请退款接口返回成功仅表示退款单已受理成功，具体的退款结果需依据退款结果通知及查询退款的返回信息为准。
func (s *Service) Refund(ctx context.Context, data RefundRequest) (RefundResult, error) {
	client, err := s.NewClient(ctx)
	if err != nil {
		return RefundResult{}, err
	}
	r := refunddomestic.RefundsApiService{Client: client}
	return s.handleRefundResult(r.Create(ctx, data.Adapter()))
}

// QueryRefund 查询单笔退款（通过商户退款单号）
// 提交退款申请后，推荐每间隔1分钟调用该接口查询一次退款状态，若超过5分钟仍是退款处理中状态，建议开始逐步衰减查询频率(比如之后间隔5分钟、10分钟、20分钟、30分钟……查询一次)。
// 退款有一定延时，零钱支付的订单退款一般5分钟内到账，银行卡支付的订单退款一般1-3个工作日到账。
// 同一商户号查询退款频率限制为300qps，如返回FREQUENCY_LIMITED频率限制报错可间隔1分钟再重试查询。
func (s *Service) QueryRefund(ctx context.Context, outTradeNo string, subMchIds ...string) (RefundResult, error) {
	client, err := s.NewClient(ctx)
	if err != nil {
		return RefundResult{}, err
	}
	r := refunddomestic.RefundsApiService{Client: client}
	req := refunddomestic.QueryByOutRefundNoRequest{
		OutRefundNo: core.String(outTradeNo),
		SubMchid:    core.String(afmt.First(subMchIds)),
	}
	return s.handleRefundResult(r.QueryByOutRefundNo(ctx, req))
}

// RefundOnAbnormal 发起异常退款
// 提交退款申请后，退款结果通知或查询退款确认状态为退款异常，可调用此接口发起异常退款处理。支持退款至用户、退款至交易商户银行账户两种处理方式。
//func (s *Service) RefundOnAbnormal(ctx context.Context, refundNo string) (RefundResult, error) {
//
//}
