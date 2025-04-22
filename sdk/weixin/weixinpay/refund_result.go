package weixinpay

import (
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/services/refunddomestic"
	"time"
)

// RefundResult 该结果仅表示退款单已受理成功，具体的退款结果需依据退款结果通知及查询退款的返回信息为准。
// refunddomestic.Refund
type RefundResult struct {
	RefundId            string                 `json:"refund_id"`             // 微信支付退款号
	OutRefundNo         string                 `json:"out_refund_no"`         // 商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	TransactionId       string                 `json:"transaction_id"`        // 微信支付交易订单号
	OutTradeNo          string                 `json:"out_trade_no"`          // 原支付交易对应的商户订单号
	Channel             refunddomestic.Channel `json:"channel"`               // 退款通道
	UserReceivedAccount string                 `json:"user_received_account"` // 取当前退款单的退款入账方，有以下几种情况： 1）退回银行卡：{银行名称}{卡类型}{卡尾号} 2）退回支付用户零钱:支付用户零钱 3）退还商户:商户基本账户商户结算银行账户 4）退回支付用户零钱通:支付用户零钱通

	SuccessAt       atype.Datetime              `json:"success_at"`
	CreatedAt       atype.Datetime              `json:"created_at"`
	Status          refunddomestic.Status       `json:"status"`
	FundsAccount    refunddomestic.FundsAccount `json:"funds_account"`
	Amount          *RefundAmount               `json:"amount"`
	PromotionDetail []RefundPromotion           `json:"promotion_detail,omitempty"`
}

func toRefundResult(resp refunddomestic.Refund, loc *time.Location) RefundResult {
	return RefundResult{
		RefundId:            types.Deref(resp.RefundId),
		OutRefundNo:         types.Deref(resp.OutRefundNo),
		TransactionId:       types.Deref(resp.TransactionId),
		OutTradeNo:          types.Deref(resp.OutTradeNo),
		Channel:             types.Deref(resp.Channel),
		UserReceivedAccount: types.Deref(resp.UserReceivedAccount),
		SuccessAt:           atype.ToDatetime2(resp.SuccessTime, loc),
		CreatedAt:           atype.ToDatetime2(resp.CreateTime, loc),
		Status:              types.Deref(resp.Status),
		FundsAccount:        types.Deref(resp.FundsAccount),
		Amount:              toRefundAmount(resp.Amount),
		PromotionDetail:     toRefundPromotions(resp.PromotionDetail),
	}
}

// RefundFromItem 多个支付账户
// refunddomestic.FundsFromItem
type RefundFromItem struct {
	Account refunddomestic.Account `json:"account"`
	Amount  atype.Money            `json:"amount"`
}

func toRefundFromItem(item refunddomestic.FundsFromItem) RefundFromItem {
	return RefundFromItem{
		Account: types.Deref(item.Account),
		Amount:  base.ToMoney(item.Amount),
	}
}
func toRefundFromItems(items []refunddomestic.FundsFromItem) []RefundFromItem {
	if len(items) == 0 {
		return nil
	}
	refunds := make([]RefundFromItem, len(items))
	for i, item := range items {
		refunds[i] = toRefundFromItem(item)
	}
	return refunds
}

// RefundAmount
// refunddomestic.Amount
type RefundAmount struct {
	Total  atype.Money `json:"total"`
	Refund atype.Money `json:"refund"`
	// 退款出资的账户类型及金额信息
	From             []RefundFromItem `json:"from"`
	PayerTotal       atype.Money      `json:"payer_total"`
	PayerRefund      atype.Money      `json:"payer_refund"`      // 退款给用户的金额，不包含所有优惠券金额
	SettlementRefund atype.Money      `json:"settlement_refund"` // 去掉非充值代金券退款金额后的退款金额，单位为分，退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额
	SettlementTotal  atype.Money      `json:"settlement_total"`  // 应结订单金额=订单金额-免充值代金券金额，应结订单金额<=订单金额，单位为分
	DiscountRefund   atype.Money      `json:"discount_refund"`   // 优惠退款金额<=退款金额，退款金额-代金券或立减优惠退款金额为现金，说明详见代金券或立减优惠，单位为分
	Currency         Currency         `json:"currency"`
}

func toRefundAmount(amount *refunddomestic.Amount) *RefundAmount {
	if amount == nil {
		return nil
	}
	return &RefundAmount{
		Total:            base.ToMoney(amount.Total),
		Refund:           base.ToMoney(amount.Refund),
		From:             toRefundFromItems(amount.From),
		PayerTotal:       base.ToMoney(amount.PayerTotal),
		PayerRefund:      base.ToMoney(amount.PayerRefund),
		SettlementRefund: base.ToMoney(amount.SettlementRefund),
		SettlementTotal:  base.ToMoney(amount.SettlementTotal),
		DiscountRefund:   base.ToMoney(amount.DiscountRefund),
		Currency:         toCurrency(amount.Currency),
	}
}

type RefundPromotion struct {
	PromotionId  string               `json:"promotion_id"` // 券或者立减优惠id
	Scope        refunddomestic.Scope `json:"scope"`
	Type         refunddomestic.Type  `json:"type"`
	Amount       atype.Money          `json:"amount"`        // 用户享受优惠的金额（优惠券面额=微信出资金额+商家出资金额+其他出资方金额 ），单位为分
	RefundAmount atype.Money          `json:"refund_amount"` // 优惠退款金额<=退款金额，退款金额-代金券或立减优惠退款金额为用户支付的现金，说明详见代金券或立减优惠，单位为分

	// 优惠商品发生退款时返回商品信息
	GoodsDetail []RefundGoodsDetail `json:"goods_detail,omitempty"`
}

func toRefundPromotion(promotion refunddomestic.Promotion) RefundPromotion {
	return RefundPromotion{
		PromotionId:  types.Deref(promotion.PromotionId),
		Scope:        types.Deref(promotion.Scope),
		Type:         types.Deref(promotion.Type),
		Amount:       base.ToMoney(promotion.Amount),
		RefundAmount: base.ToMoney(promotion.RefundAmount),
		GoodsDetail:  toGoodsDetails(promotion.GoodsDetail),
	}
}
func toRefundPromotions(promotions []refunddomestic.Promotion) []RefundPromotion {
	if len(promotions) == 0 {
		return nil
	}
	ps := make([]RefundPromotion, len(promotions))
	for i, p := range promotions {
		ps[i] = toRefundPromotion(p)
	}
	return ps
}
