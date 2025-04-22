package jsapay

import (
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/weixin/weixinpay"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"
)

// PrepayDTO 客户端直接把这个传递会给微信即可 支付结果将通过 notify url 通知
// https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_4.shtml
// 客户端获取到该dto，直接使用 WeixinJSBridge.invoke('getBrandWCPayRequest', PrepayDTO) 即可
type PrepayDTO struct {
	OutTradeNo string `json:"-"`
	PrepayId   string `json:"-"`
	AppId      string `json:"appId"`
	Timestamp  string `json:"timeStamp"`
	NonceStr   string `json:"nonceStr"`
	Package    string `json:"package"`
	SignType   string `json:"signType"`
	PaySign    string `json:"paySign"`
}

func (s *Service) newPrepayDTO(r jsapi.PrepayWithRequestPaymentResponse, outTradeNo string) PrepayDTO {
	return PrepayDTO{
		OutTradeNo: outTradeNo,
		PrepayId:   types.Deref(r.PrepayId),
		AppId:      s.Appid,
		Timestamp:  types.Deref(r.TimeStamp),
		NonceStr:   types.Deref(r.NonceStr),
		Package:    types.Deref(r.Package),
		SignType:   types.Deref(r.SignType),
		PaySign:    types.Deref(r.PaySign),
	}
}

func parsePrepayDetail(d *weixinpay.PrepayDetail) *jsapi.Detail {
	if d == nil {
		return nil
	}
	var goodsDetail []jsapi.GoodsDetail
	if len(d.GoodsDetail) > 0 {
		goodsDetail = make([]jsapi.GoodsDetail, len(d.GoodsDetail))
		for i, gd := range d.GoodsDetail {
			goodsDetail[i] = jsapi.GoodsDetail{
				MerchantGoodsId:  core.String(gd.MerchantGoodsId),
				WechatpayGoodsId: core.String(gd.WechatpayGoodsId),
				GoodsName:        core.String(gd.GoodsName),
				Quantity:         core.Int64(int64(gd.Quantity)),
				UnitPrice:        base.FromMoney(gd.UnitPrice),
			}
		}
	}
	return &jsapi.Detail{
		CostPrice:   base.FromMoney(d.CostPrice),
		InvoiceId:   core.String(d.InvoiceId),
		GoodsDetail: goodsDetail,
	}
}

func parsePrepaySceneInfo(d *weixinpay.SceneInfo) *jsapi.SceneInfo {
	if d == nil {
		return nil
	}
	var store *jsapi.StoreInfo
	if d.StoreId != "" || d.StoreName != "" || d.AreaCode != "" || d.Addr != "" {
		store = &jsapi.StoreInfo{
			Id:       core.String(d.StoreId),
			Name:     core.String(d.StoreName),
			AreaCode: core.String(d.AreaCode),
			Address:  core.String(d.Addr),
		}
	}

	return &jsapi.SceneInfo{
		PayerClientIp: core.String(d.PayerClientIp),
		DeviceId:      core.String(d.DeviceId),
		StoreInfo:     store,
	}
}
