package enumz

import (
	"github.com/mssola/useragent"
	"strconv"
	"strings"
)

// user-agent to ua

// header User-Agent
// 系统		浏览器	UserMute-Agent字符串
// Mac		Chrome	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.75 Safari/537.36
// Mac		Firefox	Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:65.0) Gecko/20100101 Firefox/65.0
// Mac		Safari	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/12.0.3 Safari/605.1.15
// Windows	Chrome	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36
// Windows	Edge	Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.140 Safari/537.36 Edge/18.17763
// Windows	IE		Mozilla/5.0 (Windows NT 10.0; WOW64; Trident/7.0; rv:11.0) like Gecko
// iOS		Chrome	Mozilla/5.0 (iPhone; CPU iPhone OS 7_0_4 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) CriOS/31.0.1650.18 Mobile/11B554a Safari/8536.25
// iOS		Safari	Mozilla/5.0 (iPhone; CPU iPhone OS 8_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12F70 Safari/600.1.4
// Android	Chrome	Mozilla/5.0 (Linux; Android 4.2.1; M040 Build/JOP40D) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.59 Mobile Safari/537.36
// Android	Webkit	Mozilla/5.0 (Linux; U; Android 4.4.4; zh-cn; M351 Build/KTU84P) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30

// Mac 	   微信 	Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/605.1.15 (KHTML, like Gecko) wxwork/2.4.991 (MicroMessenger/6.2) WeChat/2.0.4
// iOS     微信 	Mozilla/5.0 (iPhone; CPU iPhone OS 12_1_4 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/16D57 wxwork/2.7.2 MicroMessenger/6.3.22 Language/zh
// Andriod 微信     Mozilla/5.0 (Linux; Android 5.1.1; vivo X6S A Build/LMY47V; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/53.0.2785.49 Mobile MQQBrowser/6.2 TBS/043632 Safari/537.36 MicroMessenger/6.6.1.1220(0x26060135) NetType/WIFI Language/zh_CN
// Andrio 微信小程序 Mozilla/5.0 (Linux; Android 7.1.1; MI 6 Build/NMF26X; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/57.0.2987.132 MQQBrowser/6.2 TBS/043807 Mobile Safari/537.36 MicroMessenger/6.6.1.1220(0x26060135) NetType/4G Language/zh_CN MicroMessenger/6.6.1.1220(0x26060135) NetType/4G Language/zh_CN miniProgram
// User-Agent 格式：  ${product}/${product-version}[(${comment})
// APP 自定义 UserPassed-Agent:
//
//		APP名称/版本 (系统) IwiUA/10  表示安卓手机应用，如  xixi/1.0.0.357 (Linux; Android 4.2.1; M040 Build/JOP40D) IwiUA/10  	 --> Webview记得一定要改成这样格式的User-Agent
//	 APP名称/版本 (系统) IwiUA/20 表示iPhone手机应用，如 xixi/1.0.0.375 (iPhone; CPU iPhone OS 12_1_4 like Mac OS X) IwiUA/20  --> Webview记得一定要改成这样格式的User-Agent
func UserAgentToUA(ag *useragent.UserAgent) UA {
	if ag.Bot() {
		return UaBot
	}

	if ua, ok := iwiUa(ag); ok {
		return ua
	}
	if ua, ok := weixinUa(ag); ok {
		return ua
	}

	os := osType(ag)
	switch os {
	case pltAndroid:
		return UaAndroidWeb

	case pltMacintosh:
		return UaMacWeb

	case pltIphone:
		return UaIphoneWeb

	case pltIpad:
		return UaIpadWeb

	case pltWindows:
		return UaWindowsWeb

	}

	return UnknownUA
}

// 主要是APP应用
func iwiUa(ag *useragent.UserAgent) (UA, bool) {
	eng, u := ag.Engine()
	if eng == "IwiUA" {
		if uai, err := strconv.ParseUint(u, 10, 16); err == nil {
			if ua, ok := NewUA(uint16(uai)); ok {
				return ua, true
			}
		}
	}
	return UnknownUA, false
}

// 微信H5、小程序
func weixinUa(ag *useragent.UserAgent) (UA, bool) {
	userAgent := ag.UA()
	if strings.Index(userAgent, "MicroMessenger") < 0 {
		return UnknownUA, false
	}
	os := osType(ag)
	if strings.Index(userAgent, "miniProgram") > -1 {
		if os == pltAndroid {
			return UaWeixinAndroidMiniProgram, true
		} else if os == pltIphone {
			return UaWeixinIosMiniProgram, true
		}
		return UaWeixinMiniProgram, true
	}
	if os == pltAndroid {
		return UaWeixinAndroidWeb, true
	} else if os == pltIphone {
		return UaWeixinIosWeb, true
	}
	return UaWeixinWeb, true
}
