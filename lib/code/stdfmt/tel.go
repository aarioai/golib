package stdfmt

import (
	"bytes"
	"fmt"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/pkg/afmt"
	"strings"
)

const (
	TelNationCodeSeparator     = ' '
	TelAreaCodeSeparator       = '-'
	MaxNationTelCodeLength     = 5 // 含00
	MaxDistrictTelCodeLength   = 4 // 含前面0
	GlobalPhoneNumberMinLength = 3 // 110，不含国际区号、地区号
	//const GlobalPhoneNumberMaxLength = 13 //  深圳 0755-8位，共 13位，也会在不断扩容，因此不做最长判断
)

type Tel struct {
	nationCode string // 不带+和00前缀
	areaCode   string // 不带0前缀
	number     string
	readOnly   bool
}

var (
	TelZero = Tel{readOnly: true}
)

func stdTelNationCode(nationCode string) string {
	if nationCode == "" {
		return ""
	}
	if nationCode[0] == '+' {
		nationCode = nationCode[1:]
	} else if len(nationCode) > 1 && nationCode[0:2] == "00" {
		nationCode = nationCode[2:]
	}
	return nationCode
}
func stdTelAreaCode(areaCode string) string {
	if len(areaCode) > 0 && areaCode[0] == '0' {
		areaCode = areaCode[1:]
	}
	return areaCode
}

// ParseTel 格式化全球电话号码（用户手动输入），Tel 是 telephone number 官方缩写
// 如 固话 0755-2331312 和手机号 18888777999
// @return nationCode, tel
// 替换规则：
// "15000777962" =>   "15000777962",
// "15 000 777 962" =>   "15000777962",
// "（0755）12345660"=>   "0755-12345660",
// "(0755）12345660"=>  "0755-12345660",
// "(0755）- 12345660"=> "0755-12345660",
// "0755 12345660"=>    "0755-12345660",
// (0086)755 12345678  	--> +86 0755-12345678
// +86-(0)755 12345678  	--> +86 0755-12345678
// +86-(0)755 12345678  	--> +86, 0755 12345678
// +86-(755)12345678  	--> +86 0755-12345678
// +86(755)12345678		--> +86 0755-12345678
// (+86)[755]12345678  	--> +86 0755-12345678
// (0755)1234 5678  	--> 0755-12345678
// 0755)1234 5678  	--> 0755-12345678
// 美国/加拿大电话号码：11位数  NPA-NXX XXX    3位区号、3位交换代码、4位号码， 区号不用加0，如 416 表示多伦多区
// 美国地区码从2开始
// 日本电话号码：12位 080-12345678。日本国内拨打需要带0前缀，如03-123456789 国际拨打不用带0前缀， +81 3-123456789
func ParseTel(tel string, nationCodeX ...string) (Tel, error) {
	nationCode, areaCode, num, err := trimGlobalTel(tel)
	if err != nil {
		return TelZero, err
	}
	if nationCode == "" && len(nationCodeX) > 0 {
		nationCode = nationCodeX[0]
	}
	t := NewTel(nationCode, areaCode, num, false)
	if len(t.number) < GlobalPhoneNumberMinLength {
		return TelZero, ae.ErrInputWrongLength
	}
	return t, nil
}
func NewTel(nationCode, areaCode, number string, readonly bool) Tel {
	return Tel{
		nationCode: stdTelNationCode(nationCode),
		areaCode:   stdTelAreaCode(areaCode),
		number:     number,
		readOnly:   readonly,
	}
}
func NewTelWithCountry(country aenum.Country, areaCode, number string, readonly bool) Tel {
	nationCode := CountryCallingCode(country)
	return Tel{
		nationCode: stdTelNationCode(nationCode),
		areaCode:   stdTelAreaCode(areaCode),
		number:     number,
		readOnly:   readonly,
	}
}

func CountryCallingCode(country aenum.Country) string {
	switch country {
	case aenum.Canada:
		return "1"
	case aenum.NetherlandsAntilles:
		return "599"
	case aenum.WesternSahara:
		return "212"
	}
	return country.String()
}
func (t Tel) IsEmpty() bool {
	return t.number == ""
}
func (t Tel) ReadOnly() bool {
	return t.readOnly
}
func (t Tel) Reload(nationCode, areaCode, number string) {
	if t.readOnly {
		panic("reload a readonly Tel")
	}
	t.nationCode = stdTelNationCode(nationCode)
	t.areaCode = stdTelAreaCode(areaCode)
	t.number = number
}

func (t Tel) StdNationCode() string {
	return t.nationCode
}
func (t Tel) StdAreaCode() string {
	return t.areaCode
}
func (t Tel) Number() string {
	return t.number
}
func (t Tel) NationCodePad(padWithZero bool) string {
	if t.nationCode == "" {
		return ""
	}
	if padWithZero {
		return "00" + t.nationCode
	}
	return "+" + t.nationCode
}
func (t Tel) AreaCodePad() string {
	if t.areaCode == "" {
		return ""
	}
	if t.nationCode != "" {
		return t.areaCode
	}
	return "0" + t.areaCode // 不带国际码，地区码要加0前缀
}
func (t Tel) Local() string {
	if t.number == "" {
		return ""
	}

	var s strings.Builder
	s.Grow(1 + len(t.areaCode) + 1 + len(t.number))

	if t.areaCode != "" {
		s.WriteByte('0') // 这里必须要加前缀
		s.WriteString(t.areaCode)
		s.WriteByte(TelAreaCodeSeparator)
	}
	s.WriteString(t.number)
	return t.String()
}
func (t Tel) String(nationCodePadZero ...bool) string {
	if t.number == "" {
		return ""
	}
	var s strings.Builder
	s.Grow(len(t.nationCode) + 1 + len(t.areaCode) + 1 + len(t.number))
	if t.nationCode != "" {
		s.WriteString(t.NationCodePad(afmt.First(nationCodePadZero)))
		s.WriteByte(TelNationCodeSeparator)
	}
	if t.areaCode != "" {
		s.WriteString(t.AreaCodePad())
		s.WriteByte(TelAreaCodeSeparator)
	}
	s.WriteString(t.number)
	return s.String()
}

type telSegment struct {
	value    string
	enclosed bool // 被什么类型括号（开始符）括起来
}

func deepSplitTelSegments(s []byte) []telSegment {
	segs := splitTelSegments(s)

	if len(segs) < 2 {
		return segs
	}

	newSegs := make([]telSegment, 0, len(segs))
	for i := 0; i < len(segs)-1; i++ {
		next := segs[i+1]
		seg := segs[i]

		// (00)86    (0)755     188 00 334456
		if !next.enclosed && (seg.value == "00" || seg.value == "0") {
			segs[i+1].value = seg.value + next.value // 不是指针，必须要使用 segs[i+1].value =
			continue
		}
		newSegs = append(newSegs, seg)
	}

	newSegs = append(newSegs, segs[len(segs)-1])
	return newSegs
}

// 拆分不同部分
// (0086)755 12345678  	--> 0086, 755 12345678
// +86-(0)755 12345678  	--> +86, 0755  12345678
// +86-(0)755 12345678  	--> +86, 0755 12345678
// +86-(755)12345678  	--> +86, 755, 12345678
// +86(755)12345678		--> +86, 755, 12345678
// (+86)[755]12345678  	--> +86, 755, 12345678
// (0755)1234 5678  	--> 0755, 1234 5678
// 0755)1234 5678  	--> 0755, 1234 5678
func splitTelSegments(s []byte) []telSegment {
	const trim = " -" // 去掉首尾全部 - 和空格
	segs := make([]telSegment, 0)
	unenclosedStart := 0
	for i := 0; i < len(s); i++ {
		c := formatTelEncloseChar(s[i])
		// 这里如果出现 )，则表示出现缺少开始的括号
		if c == '(' || c == ')' {
			if i > unenclosedStart {
				p := string(bytes.Trim(s[unenclosedStart:i], trim))
				if p != "" {
					segs = append(segs, telSegment{
						value:    p,
						enclosed: c == ')',
					})
				}
			}
			unenclosedStart = i + 1
		}

		if c != '(' {
			continue
		}

		// 一直找到对应关闭符号
		endPos := len(s)
		for j := i + 1; j < len(s); j++ {
			char := formatTelEncloseChar(s[j])
			if char == '(' {
				return nil // 不允许括号嵌套，一方面真实场景几乎不存在。另一方面太复杂，容易解析错
			}
			if char == ')' {
				endPos = j
				break
			}
		}
		if endPos == len(s) {
			unenclosedStart = endPos
		} else {
			unenclosedStart = endPos + 1 // 最后一个 ) 不用
		}
		// 若出现 () 中间没有东西的情况，则跳过
		if endPos > i+1 {
			v := string(bytes.Trim(s[i+1:endPos], trim))
			if v != "" {
				segs = append(segs, telSegment{
					value:    v,
					enclosed: true,
				})
			}

		}
		i = endPos
	}
	if unenclosedStart < len(s) {
		v := string(bytes.Trim(s[unenclosedStart:], trim))
		if v != "" {
			segs = append(segs, telSegment{
				value:    v,
				enclosed: false,
			})
		}
	}
	return segs
}
func formatTelEncloseChar(b byte) byte {
	if b == '{' || b == '[' {
		return '('
	} else if b == '}' || b == ']' {
		return ')'
	}
	return b
}

// 替换掉多余空格
// +86-(755) 12345678  --> +86-(755)12345678
// +86 (755) 12345678  --> +86(755)12345678
// +86[755 ] 12345678    --> +86[755]123456789
// (0755) 1234 5678    --> (0755)1234 5678
func trimGlobalTelBlanks(tel string) string {
	s := make([]byte, 0, len(tel))
	s = append(s, tel[0])

	for i := 1; i < len(tel); i++ {
		c := tel[i]
		if c >= '0' && c <= '9' {
			s = append(s, c)
			continue
		}
		prevC := s[len(s)-1]
		if prevC == ' ' {
			s[len(s)-1] = c
		} else if (prevC >= '0' && prevC <= '9') || c != ' ' {
			s = append(s, c)
		}
	}
	return string(s)
}

var (
	// 国际区号是2位数的，不存在该前缀3位数国际区号
	twoDigitsNationCodes = map[string]struct{}{
		"20": {}, "27": {},
		"30": {}, "31": {}, "32": {}, "33": {}, "34": {}, "36": {}, "39": {},
		"40": {}, "41": {}, "43": {}, "44": {}, "45": {}, "46": {}, "47": {}, "48": {}, "49": {},
		"51": {}, "52": {}, "53": {}, "54": {}, "55": {}, "56": {}, "57": {}, "58": {},
		"60": {}, "61": {}, "62": {}, "63": {}, "64": {}, "65": {}, "66": {},
		"73": {}, "74": {}, "76": {}, "77": {}, "78": {}, "79": {},
		"81": {}, "82": {}, "84": {}, "86": {},
		"90": {}, "91": {}, "92": {}, "93": {}, "94": {}, "95": {}, "98": {},
	}
)

// ExtractNationCode 从加前缀nationCode，或者电话号码中，提取无前缀nation code
// @param raw
// +8675512345
// +86 75512345
// 0086-12345679
func ExtractNationCode(raw string) (string, string, bool) {
	withPrefix := false
	s := raw
	if s[0:2] == "00" {
		withPrefix = true
		s = s[2:]
	} else if s[0] == '+' {
		withPrefix = true
		s = s[1:]
	}
	if !withPrefix || len(s) == 0 {
		return "", raw, false
	}
	trim := "- " // 两端空格和横线都要去掉
	// +1 只有北美，没有 +1xx
	// +7 是俄罗斯，但是 +7x 开头有其他国家
	if s[0] == '1' || (len(s) == 1 && s[0] == '7') {
		return s[:1], strings.Trim(s[1:], trim), true
	}
	if len(s) < 2 {
		return "", raw, false
	}
	_, ok := twoDigitsNationCodes[s[:2]]
	if ok {
		return s[:2], strings.Trim(s[2:], trim), true
	}
	// 这时候再处理 +7 情况
	if s[0] == '7' {
		return "7", strings.Trim(s[1:], trim), true
	}

	if len(s) < 3 {
		return "", raw, false
	}
	return s[:3], strings.Trim(s[3:], trim), true // 其他国家都是3位数国际区号
}

// 区号是不带0开头的，不过国内拨打需要加0。国际拨打不需要加0
func parseChineseTel(tel string) (string, string) {
	if len(tel) < GlobalPhoneNumberMinLength {
		return "", tel
	}
	switch tel[0] {
	case '0':
		// 01x/02x 都是2位区号（不含开头0）
		if tel[1] == '1' || tel[1] == '2' {
			return tel[1:3], tel[3:] // 地区码是不带0的
		}
		return tel[0:4], tel[4:] // 其他都是3位区号（不含开头0）case '1':
	case '1':
		// 中国手机号段从 13 开始
		// 北京区号 10
		if tel[1] == '0' {
			return tel[0:2], tel[2:]
		}
		return "", tel
	case '2':
		return tel[0:2], tel[2:] // 一定是固话
	default:
		return tel[0:3], tel[3:] // 固话
	}
}
func parseNorthAmericaTel(tel string) (string, string) {
	// 1 开头的不是区号
	if tel == "" || tel[0] == '1' {
		return "", tel
	}

	if len(tel) < 4 {
		return "", tel
	}

	return tel[:3], tel[3:] // 北美区号3位数
}
func parseLocalTel(s, nationCode string) (string, string) {
	switch nationCode {
	case "1":
		return parseNorthAmericaTel(s)
	case "86":
		return parseChineseTel(s)
	}
	return "", s
}

// 可能含有国际区号
func parseTel(tel []byte, nationCodeX ...string) (string, string, string, error) {
	parts := splitTelPartial(tel)
	switch len(parts) {
	case 0:
		return "", "", "", ae.ErrEmptyInput
	case 1:
		s := parts[0]
		nationCode, other, ok := ExtractNationCode(s)
		if !ok {
			return afmt.First(nationCodeX), "", s, nil
		}
		areaCode, num := parseLocalTel(other, nationCode)
		return nationCode, areaCode, num, nil
	default:
		first := parts[0]
		if nationCode, other, ok := ExtractNationCode(first); ok {
			localParts := make([]string, 0, len(parts))
			if other != "" {
				localParts = append(localParts, other)
			}
			localParts = append(localParts, parts[1:]...)
			areaCode, num, err := parseSplitLocalTel(localParts, nationCode, true)
			return nationCode, areaCode, num, err
		}
		defaultNationCode := afmt.First(nationCodeX)
		areaCode, num, err := parseSplitLocalTel(parts, defaultNationCode, false)
		return defaultNationCode, areaCode, num, err
	}
}

func splitTelPartial(tel []byte) []string {
	sep := byte('-')
	if bytes.IndexByte(tel, sep) < 0 {
		sep = ' '
	}
	parts := make([]string, 0)
	prev := make([]byte, 0, len(tel))
	// 去掉空格，并拆分
	for _, t := range tel {
		if t == sep {
			if len(prev) > 0 {
				parts = append(parts, string(prev))
				prev = prev[:0] // 清空数据并重用，len(prev)=0, cap(prev)=len(tel)
			}
		} else if t != ' ' {
			prev = append(prev, t) // 移除空格
		}
	}
	if len(prev) > 0 {
		parts = append(parts, string(prev))
	}
	return parts
}

func parseSplitLocalTel(parts []string, stdNationCode string, isNationCodeInTelString bool) (string, string, error) {
	// (+86) 755 12345678
	switch len(parts) {
	case 0:
		return "", "", ae.ErrEmptyInput
	case 1:
		areaCode, num := parseLocalTel(parts[0], stdNationCode)
		return areaCode, num, nil
	default:
		first := parts[0]
		// 0 开头的，一定是地区码
		if first[0] == '0' && len(first) <= MaxNationTelCodeLength {
			districtTelCode := first[1:] // 地区码标准是不带0的，只有国内拨打才需要加0前缀
			num := strings.Join(parts[1:], "")
			return districtTelCode, num, nil
		}
		// 地区码，一般不以1开头；1开头的一般用于手机号码或其他号码
		if isNationCodeInTelString && first[0] != '1' && len(parts) == 2 && len(first) < MaxNationTelCodeLength-1 {
			districtTelCode := first
			num := parts[1]
			return districtTelCode, num, nil
		}
		s := strings.Join(parts, "")
		areaCode, num := parseLocalTel(s, stdNationCode)
		return areaCode, num, nil
	}
}

func trimGlobalTel(tel string, nationCodeX ...string) (nationCode, areaCode, num string, err error) {
	// 这里包括 （ ）【】 等，Number 会把 o O 都转为 0，而ASCII并不会这样，因此要转两次。
	newTel := ReplaceWithStdNumbers(tel)
	newTel, err = ReplaceToStdASCII(newTel, false)
	if err != nil {
		return "", "", "", fmt.Errorf("invalid tel %s:  %s", tel, err.Error())
	}
	newTel = trimGlobalTelBlanks(newTel)
	if len(newTel) < GlobalPhoneNumberMinLength {
		return "", "", "", ae.ErrInputWrongLength
	}

	segs := deepSplitTelSegments([]byte(newTel))
	switch len(segs) {
	case 1:
		return parseTel([]byte(segs[0].value), nationCode)
	case 2:
		s0 := segs[0].value
		s1 := segs[1].value
		var parts []string
		isNationCodeInTelString := false
		if nc, other, ok := ExtractNationCode(s0); ok {
			isNationCodeInTelString = true
			nationCode = nc
			if other != "" {
				num = strings.ReplaceAll(s1, " ", "")
				// (+860755) 12345678
				if other[0] == '0' && len(other) <= MaxDistrictTelCodeLength {
					areaCode = other[1:]
					return nationCode, areaCode, num, nil
				}
				// (+86755) 12345678
				if other[0] != '1' && len(other) < MaxDistrictTelCodeLength {
					areaCode = other
					return nationCode, areaCode, num, nil
				}
				// (+86189 1314 8888)132
				return "", "", "", ae.ErrInvalidInput
			} else {
				parts = splitTelPartial([]byte(s1))
			}
		} else {
			parts = splitTelPartial([]byte(s0 + " " + s1))
			nationCode = afmt.First(nationCodeX)
		}
		areaCode, num, err = parseSplitLocalTel(parts, nationCode, isNationCodeInTelString)
		return nationCode, areaCode, num, err
	case 3:
		nc, other, ok := ExtractNationCode(segs[0].value)
		if !ok || other != "" {
			return "", "", "", ae.ErrInvalidInput
		}
		nationCode = nc
		parts := make([]string, len(segs)-1)
		for i := 1; i < len(segs); i++ {
			parts[i-1] = segs[i].value
		}
		areaCode, num, err = parseSplitLocalTel(parts, nationCode, true)
		return nationCode, areaCode, num, err
	default:
		return "", "", "", ae.ErrInvalidInput
	}
}
