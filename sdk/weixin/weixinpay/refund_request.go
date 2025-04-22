package weixinpay

import (
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
)

// RefundRequest
// refunddomestic.CreateRequest
type RefundRequest struct {
	SubMchid      string                         `json:"sub_mchid"`      // 子商户的商户号，由微信支付生成并下发。服务商模式下必须传递此参数
	TransactionId string                         `json:"transaction_id"` // 原支付交易对应的微信订单号
	OutTradeNo    string                         `json:"out_trade_no"`   // 原支付交易对应的商户订单号
	OutRefundNo   string                         `json:"out_refund_no"`  // 商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	Reason        string                         `json:"reason"`         // 若商户传入，会在下发给用户的退款消息中体现退款原因
	NotifyUrl     string                         `json:"notify_url"`     // 异步接收微信支付退款结果通知的回调地址，商户平台-交易中心-退款管理-退款配置
	FundsAccount  refunddomestic.ReqFundsAccount `json:"funds_account"`  // 若传递此参数则使用对应的资金账户退款，否则默认使用未结算资金退款（仅对老资金流商户适用）  枚举值： - AVAILABLE：可用余额账户    * `AVAILABLE` - 可用余额
	Amount        RefundRequestAmount            `json:"amount"`         // 订单金额信息
	GoodsDetail   []RefundGoodsDetail            `json:"goods_detail"`   // 指定商品退款需要传此参数，其他场景无需传递
}

func (r RefundRequest) Adapter() refunddomestic.CreateRequest {
	return refunddomestic.CreateRequest{
		SubMchid:      core.String(r.SubMchid),
		TransactionId: core.String(r.TransactionId),
		OutTradeNo:    core.String(r.OutRefundNo),
		OutRefundNo:   core.String(r.OutRefundNo),
		Reason:        core.String(r.Reason),
		NotifyUrl:     core.String(r.NotifyUrl),
		FundsAccount:  types.Ref(r.FundsAccount),
		Amount:        r.Amount.adapter(),
		GoodsDetail:   toRefundGoodsDetailAdapters(r.GoodsDetail),
	}
}

type RefundRequestAmount struct {
	Refund   atype.Money      `json:"refund"`   // 退款金额，币种的最小单位，只能为整数，不能超过原订单支付金额。
	From     []RefundFromItem `json:"from"`     // 退款需要从指定账户出资时，传递此参数指定出资金额（币种的最小单位，只能为整数）。 同时指定多个账户出资退款的使用场景需要满足以下条件：1、未开通退款支出分离产品功能；2、订单属于分账订单，且分账处于待分账或分账中状态。 参数传递需要满足条件：1、基本账户可用余额出资金额与基本账户不可用余额出资金额之和等于退款金额；2、账户类型不能重复。 上述任一条件不满足将返回错误
	Total    atype.Money      `json:"total"`    // 原支付交易的订单总金额，币种的最小单位，只能为整数。
	Currency Currency         `json:"currency"` // 符合ISO 4217标准的三位字母代码，目前只支持人民币：CNY。
}

func (r *RefundRequestAmount) adapter() *refunddomestic.AmountReq {
	return &refunddomestic.AmountReq{
		Refund:   base.FromMoney(r.Refund),
		From:     nil,
		Total:    base.FromMoney(r.Total),
		Currency: core.String(r.Currency.ISO4217()),
	}
}
