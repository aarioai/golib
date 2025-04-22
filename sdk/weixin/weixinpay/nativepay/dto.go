package jsapipay

import (
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/weixin/weixinpay"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/native"
	"html/template"
)

// PrepayDTO 客户端直接把这个传递会给微信即可 支付结果将通过 notify url 通知
// https://pay.weixin.qq.com/doc/v3/merchant/4012791877
// 客户端获取到该dto，直接使用 WeixinJSBridge.invoke('getBrandWCPayRequest', PrepayDTO) 即可
type PrepayDTO struct {
	OutTradeNo string       `json:"-"`
	CodeUrl    template.URL `json:"code_url"` // 二维码地址，扫码支付
}

func (s *Service) newPrepayDTO(r native.PrepayResponse, outTradeNo string) PrepayDTO {
	return PrepayDTO{
		OutTradeNo: outTradeNo,
		CodeUrl:    template.URL(types.Deref(r.CodeUrl)),
	}
}

func parsePrepayDetail(d *weixinpay.PrepayDetail) *native.Detail {
	if d == nil {
		return nil
	}
	var goodsDetail []native.GoodsDetail
	if len(d.GoodsDetail) > 0 {
		goodsDetail = make([]native.GoodsDetail, len(d.GoodsDetail))
		for i, gd := range d.GoodsDetail {
			goodsDetail[i] = native.GoodsDetail{
				MerchantGoodsId:  core.String(gd.MerchantGoodsId),
				WechatpayGoodsId: core.String(gd.WechatpayGoodsId),
				GoodsName:        core.String(gd.GoodsName),
				Quantity:         core.Int64(int64(gd.Quantity)),
				UnitPrice:        base.FromMoney(gd.UnitPrice),
			}
		}
	}
	return &native.Detail{
		CostPrice:   base.FromMoney(d.CostPrice),
		InvoiceId:   core.String(d.InvoiceId),
		GoodsDetail: goodsDetail,
	}
}

func parsePrepaySceneInfo(d *weixinpay.SceneInfo) *native.SceneInfo {
	if d == nil {
		return nil
	}
	var store *native.StoreInfo
	if d.StoreId != "" || d.StoreName != "" || d.AreaCode != "" || d.Addr != "" {
		store = &native.StoreInfo{
			Id:       core.String(d.StoreId),
			Name:     core.String(d.StoreName),
			AreaCode: core.String(d.AreaCode),
			Address:  core.String(d.Addr),
		}
	}

	return &native.SceneInfo{
		PayerClientIp: core.String(d.PayerClientIp),
		DeviceId:      core.String(d.DeviceId),
		StoreInfo:     store,
	}
}
