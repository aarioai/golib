package weixinpay

import (
	"github.com/aarioai/airis/aa/atype"
	"time"
)

type PrepayRequest struct {
	OrderBatch uint64      `json:"order_batch"` // 批次订单号
	Total      atype.Money `json:"total"`       // 订单总额

	PayerOpenid string     `json:"payer_openid"` // 支付者openid；APP不用传
	Description string     `json:"description"`
	TimeExpire  *time.Time `json:"time_expire"` // 订单失效时间，格式为rfc3339格式
	Attach      string     `json:"attach"`      // 附加数据，怎么传过去，回调会带回来

	GoodsTag string   `json:"goods_tag"` // 商品标记，代金券或立减优惠功能的参数。
	LimitPay []string `json:"limit_pay"` // 指定支付方式
	Fapiao   bool     `json:"fapiao"`    // 传入true时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效。

	//NotifyUrl           string        `json:"notify_url"` // 有效性：1. HTTPS；2. 不允许携带查询串。
	Currency            Currency      `json:"currency"`
	Detail              *PrepayDetail `json:"detail"`
	SceneInfo           *SceneInfo    `json:"scene_info"`
	SettleProfitSharing bool          `json:"settle_profit_sharing"` // 是否指定分账
}
type PrepayDetail struct {
	CostPrice   atype.Money   `json:"cost_price"` // 1.商户侧一张小票订单可能被分多次支付，订单原价用于记录整张小票的交易金额。 2.当订单原价与支付金额不相等，则不享受优惠。 3.该字段主要用于防止同一张小票分多次支付，以享受多次优惠的情况，正常支付订单不必上传此参数。
	InvoiceId   string        `json:"invoice_id"` // 商家小票ID。
	GoodsDetail []GoodsDetail `json:"goods_detail"`
}

// SceneInfo 支付场景描述
type SceneInfo struct {
	PayerClientIp string `json:"payer_client_ip"` // 用户终端IP
	DeviceId      string `json:"device_id"`       // 商户端设备号
	StoreId       string `json:"store_id"`        // 商户门店编号
	StoreName     string `json:"store_name"`      // 商户门店名称
	AreaCode      string `json:"area_code"`       // 地区编码，详细请见微信支付提供的文档
	Addr          string `json:"addr"`            // 详细的商户门店地址
}
