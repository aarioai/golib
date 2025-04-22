package weixingzh

import (
	"context"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/golib/lib/code/coding"
	"net/url"
	"strings"
)

// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#4
//debug: true, // 开启调试模式,调用的所有api的返回值会在客户端alert出来，若要查看传入的参数，可以在pc端打开，参数信息会通过log打出，仅在pc端时才会打印。
//appId: 'wx5714806e642d1a4a', // 必填，公众号的唯一标识
//timestamp: 123, // 必填，生成签名的时间戳
//nonceStr: '', // 必填，生成签名的随机串
//signature: '',// 必填，签名
//jsApiList: ['chooseWXPay'] // 必填，需要使用的JS接口列表

type JsConfig struct {
	Debug     bool           `redis:"d" json:"debug"`
	AppID     string         `redis:"-" json:"appId"` // 这里都是直接反馈给微信的，所以需要按照微信json字段大小写规范来
	Timestamp atype.UnixTime `redis:"t" json:"timestamp"`
	NonceStr  string         `redis:"n" json:"nonceStr"`
	Signature string         `redis:"s" json:"signature"`
	JsApiList []string       `redis:"-" json:"jsApiList"`
}

func (s *Service) JsConfig(ctx context.Context, link string, apis []string, debug bool) (JsConfig,
	error) {

	nonce, ts := coding.TimeNonce()

	// 第一步拿去 js ticket
	ticket, err := s.JsTicket(ctx, "jsapi", false)

	if err != nil {
		return JsConfig{}, err
	}
	// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/JS-SDK.html#62
	// 第二步，sha1 签名  url里面所有key都要是小写
	// 参与签名的字段包括 noncestr（随机字符串）, 有效的jsapi_ticket, timestamp（时间戳）, url（当前网页的URL，不包含#及其后面部分） 。
	// link 是当前网页URL
	link, err = url.QueryUnescape(strings.Split(link, "#")[0])
	if err != nil {
		return JsConfig{}, err
	}
	params := map[string]string{
		"jsapi_ticket": ticket.Ticket,
		"noncestr":     nonce,
		"timestamp":    nonce,
		"url":          link,
	}
	signature := sha1Signature(params, "", 1024, true)
	jsconf := JsConfig{
		Debug:     debug,
		AppID:     s.appid,
		Timestamp: atype.UnixTime(ts),
		NonceStr:  nonce,
		Signature: signature,
		JsApiList: apis,
	}

	return jsconf, nil
}
