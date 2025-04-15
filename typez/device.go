package typez

import (
	"encoding/base64"
	"errors"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/google/go-querystring/query"
	"net/url"
	"strconv"
	"time"
)

type DeviceInfo struct {
	Cipher byte `json:"cipher" url:"-"` // 保留字，不需要作为URL参数

	PSID string `json:"psid" url:"psid"` // 优先级：UDID --> OAID  --> UUID  -->  pseudo_id or 浏览器finger

	UDID string `json:"udid" url:"udid"` // 【B类】Unique Device Identifier，唯一设备标识码。**相对最可靠的设备标识码**
	OAID string `json:"oaid" url:"oaid"` // 【B类】安卓广告追踪ID，等同于iPhone IDFA
	UUID string `json:"uuid" url:"uuid"` // 【C类】只是在某一时空是唯一的，当每次写在应用之后获取到的UUID都是不一样的

	Model    string `json:"model" url:"model"` // e.g. iPhone X  navigator.userAgent.deviceName
	DpWidth  uint16 `json:"dpw" url:"dpw"`     // 物理分辨率宽
	DpHeight uint16 `json:"dph" url:"dph"`     // 物理分辨率高
	DipWidth uint16 `json:"dip_w" url:"dip_w"` // 逻辑分辨率宽度  // 一定是等比例的，所以不需要高度

	UA    enumz.UA `json:"ua" url:"ua"`
	OS    string   `json:"os" url:"os"`       // ${OS类型} ${版本号}
	Agent string   `json:"agent" url:"agent"` // 浏览器就写浏览器及版本号；APP 就写APP名称及版本号
	Lang  string   `json:"lang" url:"lang"`   // 操作系统（优先）或浏览器自身语言，如 en-US 或 zh-CN 等。注意：不是浏览器accept的语言

	Info string `json:"info" url:"info"` //  扩展设备信息， 浏览器就用 user-agent

	Timestamp int64  `json:"timestamp" url:"timestamp"` // 保持混淆后每次结果不一样
	Nonce     string `json:"nonce" url:"nonce"`         // 随机码，保持每次结果不一样
}

type DeviceInfoExt struct {
	DeviceInfo
	IP string `json:"ip"` //扩展的
}

func (d DeviceInfo) Valid() bool {
	return d.PSID != ""
}

//Base64编码可用于在HTTP环境下传递较长的标识信息。例如，在Java Persistence系统Hibernate中，就采用了Base64来将一个较长的一个标识符（一般为128-bit的UUID）编码为一个字符串，用作HTTP表单和HTTP GET URL中的参数。在其他应用程序中，也常常需要把二进制数据编码为适合放在URL（包括隐藏表单域）中的形式。此时，采用Base64编码不仅比较简短，同时也具有不可读性，即所编码的数据不会被人用肉眼所直接看到。
//然而，标准的Base64并不适合直接放在URL
// Encode 将标准base64 替换  + ==> -    / ==> _ ，并且不使用 = 填充
// 浏览器通过user-agent可以默认传递很多值，就不用重复传了；浏览器不用传：info、UA
//  1. 取一个区间为[A-Z]的cipher字符，swapRange =  (cipher.ascii码 & 1) + 1
//  2. 获取时间戳 timestamp=TIMESTAMP, 8位随机字符串 nonce= rand.String() ；拼接进对象
//  3. 将对象转换为 url 传值（按key随机排序即可），注意进行 url encode  --> 不含 UA ，得到 q = "timestamp=xxx&psid=xx&udid=xxx&nonce=xxx"  如果值为空，则不参与编码
//  4. 取q长度中间位置，mid = len(q) / 2 ，ceil 取中间，保证后半段一定长于或等长于前半段
//  5. 从1位开始，每swapRange位，对m1位置右侧相应位置交换，得到q = swap(q)
//  6. 使用url友好型base64编码， 得到 b =  base64.RawURLEncoding.Encode(&b, q)
//  7. 取 b 长度中间位置， mid = len(b) / 2 ，ceil 取中间，保证后半段一定长于或等长于前半段
//  8. 从0位开始，每swapRange位，对mid位置右侧相应位置交换
//  9. 首位放上cipher字符，接上上面的字符串
// 调用频繁，越快越好，不要浪费算力

func (d DeviceInfo) Encode(cipher byte) string {
	if cipher < 'A' || cipher > 'Z' {
		cipher = 'A'
	}
	if !d.Valid() {
		return ""
	}
	d.Cipher = cipher
	swapRange := int(cipher&1) + 1
	d.Timestamp = time.Now().Unix() // 不要计算，减少解析计算成本
	d.Nonce = coding.RandAlphabets(8, d.Timestamp)
	qv, _ := query.Values(d)
	// 移除掉空值
	for k, vv := range qv {
		if len(vv) == 0 || vv[0] == "" || vv[0] == "0" {
			qv.Del(k)
		}
	}
	qs := qv.Encode()
	pp := []byte(qs)

	mid := len(pp) / 2 // ceil 取中间  --> 保证后半段一定长于或等长于前半段
	for i := 1; i < mid; i += swapRange {
		pp[i], pp[mid+i] = pp[mid+i], pp[i]
	}

	b := make([]byte, base64.RawURLEncoding.EncodedLen(len(pp)))
	base64.RawURLEncoding.Encode(b, pp)
	// 奇数位按中间对称翻转
	mid = len(b) / 2
	for i := 0; i < mid; i += swapRange {
		b[i], b[mid+i] = b[mid+i], b[i]
	}
	return string([]byte{cipher}) + string(b)
}

func DecodeDeviceInfo(s string) (*DeviceInfo, error) {
	if len(s) == 0 {
		return nil, errors.New("invalid device info")
	}
	var d DeviceInfo
	d.Cipher = s[0]
	swapRange := int(d.Cipher&1) + 1

	b := []byte(s[1:]) // 第一位是cipher
	mid := len(b) / 2
	for i := 0; i < mid; i += swapRange {
		b[i], b[mid+i] = b[mid+i], b[i]
	}

	pp := make([]byte, base64.RawURLEncoding.DecodedLen(len(b)))
	_, err := base64.RawURLEncoding.Decode(pp, b)
	if err != nil {
		return nil, err
	}

	mid = len(pp) / 2
	for i := 1; i < mid; i += swapRange {
		pp[i], pp[mid+i] = pp[mid+i], pp[i]
	}

	p, err := url.ParseQuery(string(pp))
	if err != nil {
		return nil, err
	}
	d.PSID = p.Get("psid")
	if d.PSID == "" {
		return nil, errors.New("miss psid")
	}

	d.UDID = p.Get("udid")
	d.OAID = p.Get("oaid")
	d.UUID = p.Get("uuid")
	d.Model = p.Get("model")
	ua, _ := strconv.ParseUint(p.Get("ua"), 10, 16)
	d.UA, _ = enumz.NewUA(uint16(ua))
	d.OS = p.Get("os")
	d.Agent = p.Get("agent")
	d.Lang = p.Get("lang")
	d.Info = p.Get("info")
	d.Timestamp, _ = strconv.ParseInt(p.Get("timestamp"), 10, 64)
	d.Nonce = p.Get("nonce")
	//UA:        0,
	//DpWidth:   0,
	//DpHeight:  0,
	//DipWidth:  0,

	dpw := p.Get("dpw")
	dph := p.Get("dph")
	if dpw != "" && dph != "" {
		w, _ := strconv.ParseUint(dpw, 10, 16)
		h, _ := strconv.ParseUint(dph, 10, 16)
		if w > 0 && h > 0 {
			d.DpWidth = uint16(w)
			d.DpHeight = uint16(h)
			dipW := p.Get("dip_w")
			if dipW != "" {
				dw, _ := strconv.ParseUint(dipW, 10, 16)
				if dw > 0 {
					d.DipWidth = uint16(dw)
				}
			}
		}
	}
	return &d, nil
}
