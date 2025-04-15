package enumz

import (
	"strconv"
	"strings"
)

// Base 36，所以范围是 0-36^2-1 = 0, 1295
// user-agent 是更详细的，
type UA uint16

const (
	UnknownUA UA = 0
	// 一级隐类
	UaAndroid UA = 1
	UaIos     UA = 2
	UaWindows UA = 3

	UaAndroidApp      UA = 10
	UaAndroidPhoneApp UA = 11
	UaAndroidPadApp   UA = 12
	UaAndroidPcApp    UA = 13
	UaAndroidTvApp    UA = 14

	UaAndroidWeb      UA = 15
	UaAndroidPhoneWeb UA = 16
	UaAndroidPadWeb   UA = 17
	UaAndroidPcWeb    UA = 18
	UaAndroidTvWeb    UA = 19

	UaIosApp    UA = 20
	UaIphoneApp UA = 21
	UaIpadApp   UA = 22
	UaMacApp    UA = 23
	UaIphoneWeb UA = 25
	UaIpadWeb   UA = 26
	UaMacWeb    UA = 27

	UaWindowsApp UA = 30
	UaWindowsWeb UA = 35

	UaWeixinProgram            UA = 40
	UaWeixinWeb                UA = 41
	UaWeixinAndroidWeb         UA = 42 // 在Android下的，微信公众号
	UaWeixinIosWeb             UA = 43
	UaWeixinMiniProgram        UA = 46
	UaWeixinAndroidMiniProgram UA = 47
	UaWeixinIosMiniProgram     UA = 48

	UaBot UA = 1295 // 爬虫

	//PltAlipay = 20

	// Unix    Os = 0
	// Linux   Os = 1
	// Windows Os = 2

)

var (
	UaAndroidApps        = []UA{UaAndroidPhoneApp, UaAndroidPadApp, UaAndroidPcApp, UaAndroidTvApp}
	UaIosApps            = []UA{UaIphoneApp, UaIpadApp, UaMacApp}
	UaWebs               = []UA{UaAndroidWeb, UaAndroidPhoneWeb, UaAndroidPadWeb, UaAndroidPcWeb, UaAndroidTvWeb, UaIphoneWeb, UaIpadWeb, UaMacWeb, UaWindowsWeb}
	UaWeixinWebs         = []UA{UaWeixinWeb, UaWeixinAndroidWeb, UaWeixinIosWeb}
	UaWeixinMiniPrograms = []UA{UaWeixinProgram, UaWeixinMiniProgram, UaWeixinAndroidMiniProgram, UaWeixinIosMiniProgram}
)

func NewUA(ua uint16) (UA, bool) {
	p := UA(ua)
	return p, p.Valid()
}
func (p UA) Uint16() uint16 { return uint16(p) }

// 这是 uint16，一定是正数。范围是 0-36^2-1 = 0, 1295
func (p UA) Valid() bool    { return p.Uint16() < 1295 }
func (p UA) String() string { return strconv.Itoa(int(p)) }

func ParseBase36UA(s string) (UA, bool) {
	s = strings.TrimRight(s, "_")
	n, err := strconv.ParseUint(s, 36, 16)
	if err != nil {
		return UnknownUA, false
	}
	return NewUA(uint16(n))
}

// 0-9 a-z
func (p UA) Base36(pad bool) string {
	s := strconv.FormatUint(uint64(p), 36)
	if !pad || len(s) == 2 {
		return s
	}
	return s + "_"
}
func (p UA) Name() string {
	switch p {
	case UnknownUA:
		return "未知客户端"
	case UaAndroidApp:
		return "Andriod App"
	case UaIosApp:
		return "iOS App"
	case UaWindowsApp:
		return "Windows App"
	case UaWeixinProgram:
		return "微信应用"
	case UaAndroidPadApp:
		return "Andriod Pad App"
	case UaAndroidPcApp:
		return "Andriod PC App"
	case UaAndroidTvApp:
		return "Andriod TV App"
	case UaIphoneApp:
		return "iPhone"
	case UaIpadApp:
		return "iPad"
	case UaMacApp:
		return "Mac APP"
	case UaWeixinWeb:
		return "微信H5网页"
	case UaWeixinAndroidWeb:
		return "微信H5网页(Andriod)"
	case UaWeixinIosWeb:
		return "微信H5网页(iOS)"
	case UaWeixinMiniProgram:
		return "微信小程序"
	case UaWeixinAndroidMiniProgram:
		return "微信小程序(Andriod)"
	case UaWeixinIosMiniProgram:
		return "微信小程序(iOS)"
	case UaBot:
		return "爬虫"
	}
	return "未列客户端"
}
func (p UA) Is(p2 UA) bool { return p.Uint16() == p2.Uint16() }
func (p UA) In(series []UA) bool {
	for _, plt := range series {
		if p == plt {
			return true
		}
	}
	return false
}

// 用户判断服务端是否给客户端配置cookie，网页型页面必须要设置access token cookie，否则会受到服务端权限限制。
func (p UA) IsWeb() bool { return p.In(UaWebs) || p.IsWeixinProgram() }
func (p UA) IsWeixinProgram() bool {
	return p.In(UaWeixinMiniPrograms) || p.In(UaWeixinWebs)
}
