package weixinpay

import (
	"fmt"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"
	"time"
)

type TransactionAmount struct {
	Currency      string      `json:"currency"`
	PayerCurrency string      `json:"payer_currency"`
	PayerTotal    atype.Money `json:"payer_total"` // 用户支付金额
	Total         atype.Money `json:"total"`       // 总金额
}

type Transaction struct {
	NotifyId        string                     `json:"notify_id"`
	OrderBatch      string                     `json:"order_batch"`
	Amount          TransactionAmount          `json:"amount"`
	Appid           string                     `json:"appid"`
	Attach          string                     `json:"attach"`
	BankType        string                     `json:"bank_type"`
	Mchid           string                     `json:"mchid"`
	OutTradeNo      string                     `json:"out_trade_no"`
	PayerOpenid     string                     `json:"payer_openid"`
	PromotionDetail []payments.PromotionDetail `json:"promotion_detail"`
	SuccessTime     time.Time                  `json:"success_time"`
	TradeState      TradeState                 `json:"trade_state"`
	TradeStateDesc  string                     `json:"trade_state_desc"`
	TradeType       string                     `json:"trade_type"`
	TransactionId   string                     `json:"transaction_id"`
}

func (t Transaction) Success() bool {
	return t.TradeState == TradeStateSuccess
}

// TradeState 交易状态
// https://pay.weixin.qq.com/doc/v3/merchant/4013070354
type TradeState string

const (
	TradeStateSuccess    TradeState = "SUCCESS"
	TradeStateRefund     TradeState = "REFUND"     // 转入退款
	TradeStateNotPay     TradeState = "NOTPAY"     // 未支付
	TradeStateClosed     TradeState = "CLOSED"     // 已关闭
	TradeStateRevoked    TradeState = "REVOKED"    // 已撤销（仅付款码支付会返回）
	TradeStateUserPaying TradeState = "USERPAYING" // 用户支付中（仅付款码支付会返回）
	TradeStatePayError   TradeState = "PAYERROR"   // 支付失败（仅付款码支付会返回）

	TradeStateUnknownError TradeState = "UNKNOWN_ERROR"
)

func ToTradeState(s *string) TradeState {
	if s == nil || *s == "" {
		return TradeStateUnknownError
	}
	return TradeState(*s)
}

func ToTransaction(trans payments.Transaction, notifyId string) (Transaction, error) {
	var successTime time.Time
	if trans.SuccessTime != nil && *trans.SuccessTime != "" {
		successTime, _ = time.Parse(time.RFC3339, *trans.SuccessTime)
	}

	batch, total, _, err := base.ExtractOutTradeNo(*trans.OutTradeNo)
	if err != nil {
		return Transaction{}, fmt.Errorf("extract out trade number error: %s", err.Error())
	}
	if *trans.Amount.Total != total.ToCent() {
		err = fmt.Errorf("order batch %d received payment amount cent %d(%s) is not equal to supposed %d", batch, *trans.Amount.Total, *trans.Amount.Currency, total.ToCent())
		return Transaction{}, err
	}

	t := Transaction{
		NotifyId:   notifyId,
		OrderBatch: types.FormatUint(batch),
		Amount: TransactionAmount{
			Currency:      types.Deref(trans.Amount.Currency),
			PayerCurrency: types.Deref(trans.Amount.PayerCurrency),
			PayerTotal:    base.ToMoney(trans.Amount.PayerTotal),
			Total:         base.ToMoney(trans.Amount.Total),
		},
		Appid:           types.Deref(trans.Appid),
		Attach:          types.Deref(trans.Attach),
		BankType:        types.Deref(trans.BankType),
		Mchid:           types.Deref(trans.Mchid),
		OutTradeNo:      types.Deref(trans.OutTradeNo),
		PayerOpenid:     types.Deref(trans.Payer.Openid),
		PromotionDetail: trans.PromotionDetail,
		SuccessTime:     successTime,
		TradeState:      ToTradeState(trans.TradeState),
		TradeStateDesc:  types.Deref(trans.TradeStateDesc),
		TradeType:       types.Deref(trans.TradeType),
		TransactionId:   types.Deref(trans.TransactionId),
	}
	return t, nil
}
