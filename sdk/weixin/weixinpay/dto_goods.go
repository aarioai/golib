package weixinpay

import "github.com/aarioai/airis/aa/atype"

type GoodsDetail struct {
	MerchantGoodsId  string      `json:"merchant_goods_id"`  // 由半角的大小写字母、数字、中划线、下划线中的一种或几种组成。
	WechatpayGoodsId string      `json:"wechatpay_goods_id"` // 微信支付定义的统一商品编号（没有可不传）。
	GoodsName        string      `json:"goods_name"`         // 商品的实际名称。
	Quantity         uint        `json:"quantity"`           // 用户购买的数量。
	UnitPrice        atype.Money `json:"unit_price"`         // 商品单价
}
