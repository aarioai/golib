package jsapipay

import (
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/weixin/weixinpay"
	"github.com/aarioai/golib/sdk/weixin/weixinpay/base"
	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/h5"
	"html/template"
)

// PrepayDTO 客户端直接把这个传递会给微信即可 支付结果将通过 notify url 通知
// https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_4.shtml
// 客户端获取到该dto，直接使用 WeixinJSBridge.invoke('getBrandWCPayRequest', PrepayDTO) 即可
type PrepayDTO struct {
	OutTradeNo string       `json:"-"`
	Redirect   template.URL `json:"redirect"`
}

func (s *Service) newPrepayDTO(r h5.PrepayResponse, outTradeNo string) PrepayDTO {
	return PrepayDTO{
		OutTradeNo: outTradeNo,
		Redirect:   template.URL(types.Deref(r.H5Url)),
	}
}

func parsePrepayDetail(d *weixinpay.PrepayDetail) *h5.Detail {
	if d == nil {
		return nil
	}
	var goodsDetail []h5.GoodsDetail
	if len(d.GoodsDetail) > 0 {
		goodsDetail = make([]h5.GoodsDetail, len(d.GoodsDetail))
		for i, gd := range d.GoodsDetail {
			goodsDetail[i] = h5.GoodsDetail{
				MerchantGoodsId:  core.String(gd.MerchantGoodsId),
				WechatpayGoodsId: core.String(gd.WechatpayGoodsId),
				GoodsName:        core.String(gd.GoodsName),
				Quantity:         core.Int64(int64(gd.Quantity)),
				UnitPrice:        base.FromMoney(gd.UnitPrice),
			}
		}
	}
	return &h5.Detail{
		CostPrice:   base.FromMoney(d.CostPrice),
		InvoiceId:   core.String(d.InvoiceId),
		GoodsDetail: goodsDetail,
	}
}

func parsePrepaySceneInfo(d *weixinpay.SceneInfo) *h5.SceneInfo {
	if d == nil {
		return nil
	}
	var store *h5.StoreInfo
	if d.StoreId != "" || d.StoreName != "" || d.AreaCode != "" || d.Addr != "" {
		store = &h5.StoreInfo{
			Id:       core.String(d.StoreId),
			Name:     core.String(d.StoreName),
			AreaCode: core.String(d.AreaCode),
			Address:  core.String(d.Addr),
		}
	}

	return &h5.SceneInfo{
		PayerClientIp: core.String(d.PayerClientIp),
		DeviceId:      core.String(d.DeviceId),
		StoreInfo:     store,
	}
}
