package weixinpay

import (
	"context"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"net/http"
)

type RefundStatus string

const (
	RefundSuccess  RefundStatus = "REFUND.SUCCESS"  // 退款成功
	RefundAbnormal RefundStatus = "REFUND.ABNORMAL" // 退款异常
	RefundClosed   RefundStatus = "REFUND.CLOSED"   // 退款关闭
)

// 退款结果通知
// https://pay.weixin.qq.com/doc/v3/merchant/4013070388
func (s *Service) HandleRefundNotify(ctx context.Context, r *http.Request) (RefundResult, RefundStatus, error) {
	trans := &refunddomestic.Refund{}
	handler, err := s.getRsaNotifyHandler()
	if err != nil {
		return RefundResult{}, "", err
	}

	notifyResp, err := handler.ParseNotifyRequest(ctx, r, &trans)
	if err != nil {
		return RefundResult{}, "", s.NewError("parse refund notify request error: %s", err.Error())
	}

	t := toRefundResult(*trans, s.loc)
	return t, RefundStatus(notifyResp.EventType), nil
}
