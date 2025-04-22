# 说明文档
## 方法说明
* weixingzh             微信公众号
  * GrantAccessToken      获取基础API access token (grant type: client_credential)
  * JsConfig              获取JsConfig，用户JS客户端使用SDK
  * JsTicket              获取JsTicket
  * BatchUserinfo         批量获取微信用户资料
* weixinwebauth         微信网页认证
  * Auth                  通过code换取网页登录access_token (grant type: authorization_code)。注意：这里不是基础API access token
  * UserInfo              获取用户信息
  * Code2UserInfo         = Auth + UserInfo
* weixinjsa              微信小程序
* weixinpay             微信支付
  * apppay|h5pay|jsapay|nativepay    APP支付|公众号H5支付|小程序支付|QR扫码支付
    * Prepay                下单，返回不同平台对应的 PrepayDTO 
    * CloseOrder            关闭订单
    * QueryUnpaidOrder      未支付成功订单查询
    * QueryPaidOrder        支付成功订单查询
    * PayService
      * Refund                退款申请
      * QueryRefund           查询单笔退款（通过商户退款单号）
      * HandlePayNotify       notify_url 支付事件处理
      * HandleRefundNotify
 * weixinpartner            服务商
 * weixintransfer           商家转账
   * Transfer                 发起商家转账
   * Query                    通过商家批次单号查询批次单
   * QueryByWeixinBatch       通过微信批次单号查询批次单
   * QueryDetail              通过商家明细单号查询明细单
   * QueryDetailByWeixinBatch 通过微信明细单号查询明细单


## 微信H5支付示意图
```
 [微信H5] 
   +----- 1.1 请求后端下单 -------------------------------->  [后端]
   | <--- 1.2 后端返回 PrepayDTO（含prepay_id） ------------
   |                                                                              
   +----- 2.1 通过 WeixinJSBridge(PrepayDTO) 下单 -------->  [微信服务器] 
   | <--- 2.2 微信返回 h5_url，跳转到该页面 -----------------
   |
   +----- 3.1 用户在 h5_url 页，发起支付 ---------------------  [微信服务器] -- 4.0 通过微信支付后台配置的 notify_url 通知 ---> [后端]
```

## 微信小程序支付示意图
```
 [微信小程序] 
   +----- 1.1 请求后端下单 -------------------------------->  [后端]
   | <--- 1.2 后端返回 PrepayDTO（含 prepay_id） ------------
   |                                                                              
   +----- 2.1 wx.requestPayment(PrepayDTO)，发起支付 ----->   [微信服务器] -- 3.0 通过微信支付后台配置的 notify_url 通知 ---> [后端]
```


## 微信APP支付示意图
```
[微信APP] <---- 1.2 后端返回 PrepayDTO ------- [后端]
    |      ---------------1.1 -------------> 
    |
    +-- 2. wx.sendReq(JsPayDTO) 发起微信支付 ---------------------------> [微信服务器]
                                                                                   |
                                           [后端] <----- 3. notify url 通知 --------+
```

## 微信JSPay
```
var (
	wxpayConfigKey = weixinpay.ConfigKey{
		// 微信公众号部分
		AppidKeyName: "wxgzh.appid",

		// 微信支付部分
		MchIdKeyName:         "wxpay.mchid",
		MchCertSerialKeyName: "wxpay.mch_cert_serial",
		MchApiKeyKeyName:     "wxpay.mch_api_key",
		PermsDirKeyName:      "wxpay.perms_dir",
		NotifyUrlKeyName:     "wxpay.notify_url",
	}
)

// 微信JSPay接口
func (s *Service) PrepareWeixinJspay(ctx context.Context, openid string, orderBatch string, total atype.Money, retryAt atype.Datetime) (weixinpay.PrepayDTO, *ae.Error) {
	wxpay := weixinpay.New(s.app).WithConfigKey(wxpayConfigKey)
	data, err := wxpay.PrepareWeixinJspay(ctx, openid, orderBatch, total, retryAt)
	if err != nil {
		return weixinpay.PrepayDTO{}, ae.NewError(err)
	}
	return data, nil
}
```