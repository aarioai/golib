package weixintransfer

import (
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"
	"time"
)

type DetailStatus string

const (
	DetailSuccess    QueryStatus = "SUCCESS"    // 转账成功
	DetailFailed     QueryStatus = "FAIL"       // 转账失败，需要确认失败原因后，再决定是否重新发起对该笔明细单的转账（并非整个转账批次单）
	DetailProcessing QueryStatus = "PROCESSING" // 正在处理中，转账结果尚未明确
)

// DetailRequest
// transferbatch.GetTransferDetailByOutNoRequest
type DetailRequest struct {
	OutBatchNo  string `json:"out_batch_no"`  // 微信支付批次单号，微信商家转账系统返回的唯一标识
	OutDetailNo string `json:"out_detail_no"` // 微信支付系统内部区分转账批次单下不同转账明细单的唯一标识
}

func (r DetailRequest) adapter() transferbatch.GetTransferDetailByOutNoRequest {
	return transferbatch.GetTransferDetailByOutNoRequest{
		OutBatchNo:  core.String(r.OutBatchNo),
		OutDetailNo: core.String(r.OutDetailNo),
	}
}

// DetailRequestByWeixinBatch
// transferbatch.GetTransferDetailByNoRequest
type DetailRequestByWeixinBatch struct {
	BatchId  string `json:"batch_id"`  // 微信支付批次单号，微信商家转账系统返回的唯一标识
	DetailId string `json:"detail_id"` // 微信支付系统内部区分转账批次单下不同转账明细单的唯一标识
}

func (r DetailRequestByWeixinBatch) adapter() transferbatch.GetTransferDetailByNoRequest {
	return transferbatch.GetTransferDetailByNoRequest{
		BatchId:  core.String(r.BatchId),
		DetailId: core.String(r.DetailId),
	}
}

type DetailResult struct {
	Mchid          string                       `json:"mchid"`         // 微信支付分配的商户号，此处为服务商商户号
	OutBatchNo     string                       `json:"out_batch_no"`  // 商户系统内部的商家批次单号，在商户系统内部唯一
	BatchId        string                       `json:"batch_id"`      // 微信支付批次单号，微信商家转账系统返回的唯一标识
	Appid          string                       `json:"appid"`         // 微信分配的商户公众账号ID。特约商户授权类型为INFORMATION_AUTHORIZATION_TYPE和INFORMATION_AND_FUND_AUTHORIZATION_TYPE时对应的是特约商户的appid，特约商户授权类型为FUND_AUTHORIZATION_TYPE时为服务商的appid
	OutDetailNo    string                       `json:"out_detail_no"` // 商户系统内部区分转账批次单下不同转账明细单的唯一标识
	DetailId       string                       `json:"detail_id"`     // 微信支付系统内部区分转账批次单下不同转账明细单的唯一标识
	DetailStatus   DetailStatus                 `json:"detail_status"`
	TransferAmount atype.Money                  `json:"transfer_amount"` // 转账金额
	TransferRemark string                       `json:"transfer_remark"` // 单条转账备注（微信用户会收到该备注），UTF8编码，最多允许32个字符
	FailReason     transferbatch.FailReasonType `json:"fail_reason"`
	Openid         string                       `json:"openid"`         // 收款用户openid。如果转账特约商户授权类型是INFORMATION_AUTHORIZATION_TYPE，对应的是特约商户公众号下的openid;如果转账特约商户授权类型是FUND_AUTHORIZATION_TYPE，对应的是服务商商户公众号下的openid。
	Username       string                       `json:"username"`       // 收款方姓名。采用标准RSA算法，公钥由微信侧提供
	InitiatedAt    atype.Datetime               `json:"initiated_time"` // 转账发起的时间
	UpdatedAt      atype.Datetime               `json:"updated_at"`     // 明细最后一次状态变更的时间
}

// toDetailResult
// transferbatch.TransferDetailEntity
func toDetailResult(t transferbatch.TransferDetailEntity, loc *time.Location) DetailResult {
	return DetailResult{
		Mchid:          types.Deref(t.Mchid),
		OutBatchNo:     types.Deref(t.OutBatchNo),
		BatchId:        types.Deref(t.BatchId),
		Appid:          types.Deref(t.Appid),
		OutDetailNo:    types.Deref(t.OutDetailNo),
		DetailId:       types.Deref(t.DetailId),
		DetailStatus:   DetailStatus(types.Deref(t.DetailStatus)),
		TransferAmount: base.ToMoney(t.TransferAmount),
		TransferRemark: types.Deref(t.TransferRemark),
		FailReason:     types.Deref(t.FailReason),
		Openid:         types.Deref(t.Openid),
		Username:       types.Deref(t.UserName),
		InitiatedAt:    atype.ToDatetime2(t.InitiateTime, loc),
		UpdatedAt:      atype.ToDatetime2(t.UpdateTime, loc),
	}
}
