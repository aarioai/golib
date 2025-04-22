package weixintransfer

import (
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/transferbatch"
	"time"
)

type TransferType string

const (
	TransferByAPI TransferType = "API"
	TransferByWeb TransferType = "WEB"
)

type TransferStatus string

const (
	TransferInit       TransferStatus = "INIT"       // 初始态。 系统转账校验中
	TransferAccepted   TransferStatus = "ACCEPTED"   // 批次已受理，若发起批量转账的30分钟后，转账批次单仍处于该状态，可能原因是商户账户余额不足等。商户可查询账户资金流水，若该笔转账批次单的扣款已经发生，则表示批次已经进入转账中，请再次查单确认
	TransferProcessing TransferStatus = "PROCESSING" // 转账中。已开始处理批次内的转账明细单
	TransferFinished   TransferStatus = "FINISHED"   // 批次内的所有转账明细单都已处理完成
	TransferClosed     TransferStatus = "CLOSED"     // 可查询具体的批次关闭原因确认
	TransferWaitPay    TransferStatus = "WAIT_PAY"
)

// TransferRequest
// transferbatch.InitiateBatchTransferRequest
type TransferRequest struct {
	OutBatchNo  string      `json:"out_batch_no"` // 商户系统内部的商家批次单号，要求此参数只能由数字、大小写字母组成，在商户系统内部唯一
	BatchName   string      `json:"batch_name"`   // 该笔批量转账的名称
	BatchRemark string      `json:"batch_remark"` // 转账说明，UTF8编码，最多允许32个字符
	TotalAmount atype.Money `json:"total_amount"` // 转账金额。转账总金额必须与批次内所有明细转账金额之和保持一致，否则无法发起转账操作
	TotalNum    int64       `json:"total_num"`    // 一个转账批次单最多发起一千笔转账。转账总笔数必须与批次内所有明细之和保持一致，否则无法发起转账操作

	TransferDetailList []TransferDetailInput `json:"transfer_detail_list"` // 发起批量转账的明细列表，最多一千笔
	TransferSceneId    string                `json:"transfer_scene_id"`    // 该批次转账使用的转账场景，如不填写则使用商家的默认场景，如无默认场景可为空，可前往“商家转账到零钱-前往功能”中申请。 如：1001-现金营销
	NotifyUrl          string                `json:"notify_url"`           // 商户接收批次结果通知的URL，必须支持https，且只能是直接可访问的URL，不允许携带查询参数

}
type TransferDetailInput struct {
	OutDetailNo    string      `json:"out_detail_no"`   // 商户系统内部区分转账批次单下不同转账明细单的唯一标识，要求此参数只能由数字、大小写字母组成
	TransferAmount atype.Money `json:"transfer_amount"` // 转账金额
	TransferRemark string      `json:"transfer_remark"` // 单条转账备注（微信用户会收到该备注），UTF8编码，最多允许32个字符
	Openid         string      `json:"openid"`          // 商户appid下，某用户的openid
	UserName       string      `json:"user_name"`       // 收款人姓名。明细转账金额>=2000元必须要传收款人姓名
}

func (t TransferDetailInput) adapter() transferbatch.TransferDetailInput {
	// 收款方真实姓名。明细转账金额<0.3元时，不允许填写收款用户姓名
	// 明细转账金额 >= 2,000元时，该笔明细必须填写收款用户姓名
	// 同一批次转账明细中的姓名字段传入规则需保持一致，也即全部填写、或全部不填写
	// 若商户传入收款用户姓名，微信支付会校验用户openID与姓名是否一致，并提供电子回单
	if t.TransferAmount < 3*atype.Dime {
		t.UserName = ""
	}
	return transferbatch.TransferDetailInput{
		OutDetailNo:    core.String(t.OutDetailNo),
		TransferAmount: base.FromMoney(t.TransferAmount),
		TransferRemark: core.String(t.TransferRemark),
		Openid:         core.String(t.Openid),
		UserName:       core.String(t.UserName),
	}
}

func (t TransferRequest) adapter(appid string) transferbatch.InitiateBatchTransferRequest {
	details := make([]transferbatch.TransferDetailInput, len(t.TransferDetailList))
	for i, detail := range t.TransferDetailList {
		details[i] = detail.adapter()
	}
	return transferbatch.InitiateBatchTransferRequest{
		Appid:              core.String(appid),
		OutBatchNo:         core.String(t.OutBatchNo),
		BatchName:          core.String(t.BatchName),
		BatchRemark:        core.String(t.BatchRemark),
		TotalAmount:        base.FromMoney(t.TotalAmount),
		TotalNum:           core.Int64(t.TotalNum),
		TransferDetailList: details,
		TransferSceneId:    core.String(t.TransferSceneId),
		NotifyUrl:          core.String(t.NotifyUrl),
	}
}

type TransferResult struct {
	OutBatchNo  string         `json:"out_batch_no"`
	BatchId     string         `json:"batch_id"` // 微信批次单号，微信商家转账系统返回的唯一标识
	CreatedAt   atype.Datetime `json:"created_at"`
	BatchStatus TransferStatus `json:"batch_status"`
}

func toTransferResponse(t transferbatch.InitiateBatchTransferResponse, loc *time.Location) TransferResult {
	return TransferResult{
		OutBatchNo:  types.Deref(t.OutBatchNo),
		BatchId:     types.Deref(t.BatchId),
		CreatedAt:   atype.ToDatetime2(t.CreateTime, loc),
		BatchStatus: TransferStatus(types.Deref(t.BatchStatus)),
	}
}
