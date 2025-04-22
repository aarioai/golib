package weixinpay

// Currency 符合ISO 4217标准的三位字母代码，目前只支持人民币：CNY。
type Currency string

const (
	CNY Currency = "CNY"
)

func (c Currency) ISO4217() string {
	if c == "" {
		return string(c)
	}
	return string(c)
}
func toCurrency(s *string) Currency {
	if s == nil || *s == "" {
		return CNY
	}
	return Currency(*s)
}

type AppTypeRaw string

const (
	AtrH5     AppTypeRaw = "h5"     // 公众号H5
	AtrApp    AppTypeRaw = "app"    // APP调取微信
	AtrJsa    AppTypeRaw = "jsapi"  // 小程序
	AtrNative AppTypeRaw = "native" // 二维码支付
)

type AppType uint8

const (
	AtUnknown AppType = 0
	AtH5      AppType = 1 // 公众号H5
	AtApp     AppType = 2 // APP调取微信
	AtJsa     AppType = 3 // 小程序
	AtNative  AppType = 4 // 二维码支付
)

func (t AppTypeRaw) IsValid() bool {
	return t == AtrNative || t == AtrJsa || t == AtrH5 || t == AtrApp
}

func (t AppTypeRaw) Type() AppType {
	switch t {
	case AtrH5:
		return AtH5
	case AtrApp:
		return AtApp
	case AtrJsa:
		return AtJsa
	case AtrNative:
		return AtNative
	default:
		return AtUnknown
	}
}
