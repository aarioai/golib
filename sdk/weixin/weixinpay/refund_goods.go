package weixinpay

import (
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

type RefundGoodsDetail struct {
	MerchantGoodsId  string      `json:"merchant_goods_id"`  // 由半角的大小写字母、数字、中划线、下划线中的一种或几种组成
	WechatpayGoodsId string      `json:"wechatpay_goods_id"` // 微信支付定义的统一商品编号（没有可不传）
	GoodsName        string      `json:"goods_name"`         // 商品的实际名称
	UnitPrice        atype.Money `json:"unit_price"`
	RefundAmount     atype.Money `json:"refund_amount"`
	RefundQuantity   uint        `json:"refund_quantity"`
}

func (d RefundGoodsDetail) adapter() refunddomestic.GoodsDetail {
	return refunddomestic.GoodsDetail{
		MerchantGoodsId:  core.String(d.MerchantGoodsId),
		WechatpayGoodsId: core.String(d.WechatpayGoodsId),
		GoodsName:        core.String(d.GoodsName),
		UnitPrice:        base.FromMoney(d.UnitPrice),
		RefundAmount:     base.FromMoney(d.RefundAmount),
		RefundQuantity:   types.Ref(int64(d.RefundQuantity)),
	}
}
func toRefundGoodsDetailAdapters(ds []RefundGoodsDetail) []refunddomestic.GoodsDetail {
	if len(ds) == 0 {
		return nil
	}
	res := make([]refunddomestic.GoodsDetail, len(ds))
	for i, d := range ds {
		res[i] = d.adapter()
	}
	return res
}
func toRefundGoodsDetail(detail refunddomestic.GoodsDetail) RefundGoodsDetail {
	return RefundGoodsDetail{
		MerchantGoodsId:  types.Deref(detail.MerchantGoodsId),
		WechatpayGoodsId: types.Deref(detail.WechatpayGoodsId),
		GoodsName:        types.Deref(detail.GoodsName),
		UnitPrice:        base.ToMoney(detail.UnitPrice),
		RefundAmount:     base.ToMoney(detail.RefundAmount),
		RefundQuantity:   uint(types.Deref(detail.RefundQuantity)),
	}
}
func toGoodsDetails(details []refunddomestic.GoodsDetail) []RefundGoodsDetail {
	if len(details) == 0 {
		return nil
	}
	gds := make([]RefundGoodsDetail, 0, len(details))
	for _, detail := range details {
		gds = append(gds, toRefundGoodsDetail(detail))
	}
	return gds
}
