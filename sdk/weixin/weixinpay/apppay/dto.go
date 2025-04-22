package apppay

import (
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/weixin/weixinpay"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/app"
)

// PrepayDTO 客户端直接把这个传递会给微信即可 支付结果将通过 notify url 通知
// https://pay.weixin.qq.com/doc/v3/merchant/4013070351
// 客户端获取到该dto，直接使用 sendReq(PrepayDTO) 即可
type PrepayDTO struct {
	OutTradeNo string `json:"-"`

	AppId        string `json:"appId"`
	PartnerId    string `json:"partnerId"` // 商户 mchid
	PrepayId     string `json:"prepayId"`
	PackageValue string `json:"packageValue"` // Sign=WXPay
	NonceStr     string `json:"nonceStr"`
	Timestamp    string `json:"timeStamp"` //  string类型
	Sign         string `json:"sign"`      // 签名，使用字段appId、timeStamp、nonceStr、prepayId以及商户API证书私钥生成的RSA签名值
}

func (s *Service) newPrepayDTO(r app.PrepayWithRequestPaymentResponse, outTradeNo string) PrepayDTO {
	return PrepayDTO{
		OutTradeNo:   outTradeNo,
		AppId:        s.Appid,
		PartnerId:    types.Deref(r.PartnerId),
		PrepayId:     types.Deref(r.PrepayId),
		PackageValue: types.Deref(r.Package),
		NonceStr:     types.Deref(r.NonceStr),
		Timestamp:    types.Deref(r.TimeStamp),
		Sign:         types.Deref(r.Sign),
	}
}

func parsePrepayDetail(d *weixinpay.PrepayDetail) *app.Detail {
	if d == nil {
		return nil
	}
	var goodsDetail []app.GoodsDetail
	if len(d.GoodsDetail) > 0 {
		goodsDetail = make([]app.GoodsDetail, len(d.GoodsDetail))
		for i, gd := range d.GoodsDetail {
			goodsDetail[i] = app.GoodsDetail{
				MerchantGoodsId:  core.String(gd.MerchantGoodsId),
				WechatpayGoodsId: core.String(gd.WechatpayGoodsId),
				GoodsName:        core.String(gd.GoodsName),
				Quantity:         core.Int64(int64(gd.Quantity)),
				UnitPrice:        base.FromMoney(gd.UnitPrice),
			}
		}
	}
	return &app.Detail{
		CostPrice:   base.FromMoney(d.CostPrice),
		InvoiceId:   core.String(d.InvoiceId),
		GoodsDetail: goodsDetail,
	}
}

func parsePrepaySceneInfo(d *weixinpay.SceneInfo) *app.SceneInfo {
	if d == nil {
		return nil
	}
	var store *app.StoreInfo
	if d.StoreId != "" || d.StoreName != "" || d.AreaCode != "" || d.Addr != "" {
		store = &app.StoreInfo{
			Id:       core.String(d.StoreId),
			Name:     core.String(d.StoreName),
			AreaCode: core.String(d.AreaCode),
			Address:  core.String(d.Addr),
		}
	}

	return &app.SceneInfo{
		PayerClientIp: core.String(d.PayerClientIp),
		DeviceId:      core.String(d.DeviceId),
		StoreInfo:     store,
	}
}
