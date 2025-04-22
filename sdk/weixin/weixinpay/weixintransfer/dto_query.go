package weixintransfer

import (
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"
	"time"
)

type QueryStatus string

const (
	QueryWaitPay    QueryStatus = "WAIT_PAY" // 待商户确认, 符合免密条件时, 系统会自动扭转为转账中
	QueryAllSuccess QueryStatus = "ALL"      // 同时查询转账成功、失败和待确认的明细单
	QuerySuccess    QueryStatus = "SUCCESS"  // 转账成功
	QueryFailed     QueryStatus = "FAIL"     // 转账失败，需要确认失败原因后，再决定是否重新发起对该笔明细单的转账（并非整个转账批次单）
)

// QueryRequest
// transferbatch.GetTransferBatchByOutNoRequest
type QueryRequest struct {
	OutBatchNo      string      `json:"out_batch_no"`      // 商户系统内部的商家批次单号，在商户系统内部唯一
	NeedQueryDetail bool        `json:"need_query_detail"` // 商户可选择是否查询指定状态的转账明细单，当转账批次单状态为“FINISHED”（已完成）时，才会返回满足条件的转账明细单
	Offset          int64       `json:"offset"`            // 该次请求资源（转账明细单）的起始位置，从0开始，默认值为0
	Limit           int64       `json:"limit"`             // 该次请求可返回的最大资源（转账明细单）条数，最小20条，最大100条，不传则默认20条。不足20条按实际条数返回
	DetailStatus    QueryStatus `json:"detail_status"`
}

func (q QueryRequest) adapter() transferbatch.GetTransferBatchByOutNoRequest {
	return transferbatch.GetTransferBatchByOutNoRequest{
		OutBatchNo:      core.String(q.OutBatchNo),
		NeedQueryDetail: core.Bool(q.NeedQueryDetail),
		Offset:          core.Int64(q.Offset),
		Limit:           core.Int64(q.Limit),
		DetailStatus:    core.String(string(q.DetailStatus)),
	}
}

// QueryRequestByWeixinBatch
// transferbatch.GetTransferBatchByNoRequest
type QueryRequestByWeixinBatch struct {
	BatchId         string      `json:"batch_id"`
	NeedQueryDetail bool        `json:"need_query_detail"` // 商户可选择是否查询指定状态的转账明细单，当转账批次单状态为“FINISHED”（已完成）时，才会返回满足条件的转账明细单
	Offset          int64       `json:"offset"`            // 该次请求资源的起始位置。返回的明细是按照设置的明细条数进行分页展示的，一次查询可能无法返回所有明细，我们使用该参数标识查询开始位置，默认值为0
	Limit           int64       `json:"limit"`             // 该次请求可返回的最大明细条数，最小20条，最大100条，不传则默认20条。不足20条按实际条数返回
	DetailStatus    QueryStatus `json:"detail_status"`
}

func (q QueryRequestByWeixinBatch) adapter() transferbatch.GetTransferBatchByNoRequest {
	return transferbatch.GetTransferBatchByNoRequest{
		BatchId:         core.String(q.BatchId),
		NeedQueryDetail: core.Bool(q.NeedQueryDetail),
		Offset:          core.Int64(q.Offset),
		Limit:           core.Int64(q.Limit),
		DetailStatus:    core.String(string(q.DetailStatus)),
	}
}

// QueryResult
// transferbatch.TransferBatchEntity
type QueryResult struct {
	Mchid           string                        `json:"mchid"`
	OutBatchNo      string                        `json:"out_batch_no"` // 商户系统内部的商家批次单号，在商户系统内部唯一
	BatchId         string                        `json:"batch_id"`     // 微信批次单号，微信商家转账系统返回的唯一标识
	Appid           string                        `json:"appid"`        // 申请商户号的appid或商户号绑定的appid（企业号corpid即为此appid）
	BatchStatus     TransferStatus                `json:"batch_status"`
	BatchType       TransferType                  `json:"batch_type"`
	BatchName       string                        `json:"batch_name"`   // 该笔批量转账的名称
	BatchRemark     string                        `json:"batch_remark"` // 转账说明，UTF8编码，最多允许32个字符
	CloseReason     transferbatch.CloseReasonType `json:"close_reason"` // 如果批次单状态为“CLOSED”（已关闭），则有关闭原因
	TotalAmount     atype.Money                   `json:"total_amount"` // 转账金额
	TotalNum        atype.Money                   `json:"total_num"`    // 一个转账批次单最多发起三千笔转账
	CreatedAt       atype.Datetime                `json:"created_at"`
	UpdatedAt       atype.Datetime                `json:"updated_at"`
	SuccessAmount   atype.Money                   `json:"success_amount"`    // 转账成功的金额。当批次状态为“PROCESSING”（转账中）时，转账成功金额随时可能变化
	SuccessNum      int64                         `json:"success_num"`       // 转账成功的笔数。当批次状态为“PROCESSING”（转账中）时，转账成功笔数随时可能变化
	FailAmount      atype.Money                   `json:"fail_amount"`       // 转账失败的金额
	FailNum         int64                         `json:"fail_num"`          // 转账失败的笔数
	TransferSceneId string                        `json:"transfer_scene_id"` // 指定的转账场景ID

	TransferDetailList []TransferDetailCompact `json:"transfer_detail_list,omitempty"` // 当批次状态为“FINISHED”（已完成），且成功查询到转账明细单时返回。包括微信明细单号、明细状态信息
}
type TransferDetailCompact struct {
	DetailId     string         `json:"detail_id"`     // 微信支付系统内部区分转账批次单下不同转账明细单的唯一标识
	OutDetailNo  string         `json:"out_detail_no"` // 商户系统内部区分转账批次单下不同转账明细单的唯一标识
	DetailStatus TransferStatus `json:"detail_status"`
}

func toQueryResult(t transferbatch.TransferBatchEntity, loc *time.Location) QueryResult {
	var list []TransferDetailCompact
	if len(t.TransferDetailList) > 0 {
		list = make([]TransferDetailCompact, len(t.TransferDetailList))
		for i, v := range t.TransferDetailList {
			list[i] = TransferDetailCompact{
				DetailId:     types.Deref(v.DetailId),
				OutDetailNo:  types.Deref(v.OutDetailNo),
				DetailStatus: TransferStatus(types.Deref(v.DetailStatus)),
			}
		}
	}
	tb := t.TransferBatch
	if tb == nil {
		return QueryResult{
			TransferDetailList: list,
		}
	}

	return QueryResult{
		Mchid:              types.Deref(tb.Mchid),
		OutBatchNo:         types.Deref(tb.OutBatchNo),
		BatchId:            types.Deref(tb.BatchId),
		Appid:              types.Deref(tb.Appid),
		BatchStatus:        TransferStatus(types.Deref(tb.BatchStatus)),
		BatchType:          TransferType(types.Deref(tb.BatchType)),
		BatchName:          types.Deref(tb.BatchName),
		BatchRemark:        types.Deref(tb.BatchRemark),
		CloseReason:        types.Deref(tb.CloseReason),
		TotalAmount:        base.ToMoney(tb.TotalAmount),
		TotalNum:           base.ToMoney(tb.TotalNum),
		CreatedAt:          atype.ToDatetime2(tb.CreateTime, loc),
		UpdatedAt:          atype.ToDatetime2(tb.UpdateTime, loc),
		SuccessAmount:      base.ToMoney(tb.SuccessAmount),
		SuccessNum:         types.Deref(tb.SuccessNum),
		FailAmount:         base.ToMoney(tb.FailAmount),
		FailNum:            types.Deref(tb.FailNum),
		TransferSceneId:    types.Deref(tb.TransferSceneId),
		TransferDetailList: list,
	}
}
