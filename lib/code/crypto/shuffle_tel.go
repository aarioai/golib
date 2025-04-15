package crypto

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"log"
	"math/rand/v2"
	"strings"

	"github.com/aarioai/airis/aa/ae"
)

// @see _说明.md
const (
	TelCipherSeparator = byte(' ') // 除了分隔符以外其他符号，都可以
	TelKeyMinLen       = 60        // 最短要求60个。 stdfmt.SafeFastEncryptCodes 是89个，完全可以满足 telPrefixes 范围
	TelMinLength       = 6

	TelCipherSafeLen                  = 13 // 含国际码
	TelNationCodeMaxLen               = 3  // 2位64进制码 + 一个分隔符
	TelCipherMinLen                   = 6  // 不含国际区号 @see 说明文档
	TelCipherMaxLen                   = 9
	TelCipherScatterMinLen            = 7
	TelCipherScatterMaxLen            = 10
	PhoneNumberCipherMinLength        = 7 // 手机号密文长度范围（不含国际码）
	PhoneNumberCipherMaxLength        = 8
	PhoneNumberCipherScatterMinLength = 8 // 分散模式手机号密文长度范围（不含国际码）
	PhoneNumberCipherScatterMaxLength = 9

	telShiftSegmentLen        = 10 // 0-9
	telShiftSegmentItemMinLen = 6  // 每节至少6个数字，否则不利于混淆
)

// TelCipher CHAR(9)
type TelCipher string

type TelCipherText struct {
	Ciphertext string
	Tel        string
	NationCode string
	AreaCode   string
}

// ValidateTelEncryptKeys 一般启动的时候，验证key是否有错误。启动的时候执行，因此不必在意性能消耗
// 支持 除分隔符（crypto.TelCipherSeparator）外，其他所有字符。但是要求避免把types.Base64Digits全部字符都包括在内
func ValidateTelEncryptKeys[T string | []byte](keys ...T) {
	// 基于FastEncrypt
	coding.ValidateShuffleEncryptKeys(keys...)

	b64d := []byte(types.Base64Digits)
	for _, item := range keys {
		key := []byte(item)
		// 不能包含分隔符
		if bytes.Contains(key, []byte{TelCipherSeparator}) {
			panic(fmt.Sprintf("tel encrypt key %s is invalid", string(key)))
		}
		// 使用了base64转化数字， 应当包括所有 types.Base64Digits 字符
	loop:
		for _, c := range b64d {
			for _, k := range key {
				if k == c {
					continue loop
				}
			}
			panic("tel encrypt key miss base64 digit char: " + string(c))
		}

	}
}

// ShuffleEncryptTel 对电话号码加密。Tel 是 telephone number 官方缩写
// 电话区号+最后1位都是明文；手机号：前三位+最后一位是明文
// @test 保证不修改 key
func ShuffleEncryptTel(tel string, key []byte, scatter bool) (TelCipher, stdfmt.Tel, error) {
	ValidateKeyLen(key, TelKeyMinLen)
	if tel == "" {
		return "", stdfmt.TelZero, nil
	}
	t, err := stdfmt.ParseTel(tel)
	if err != nil {
		return "", stdfmt.TelZero, err
	}
	num := t.Number()
	n := len(num)
	// 少于6个字符，不加密。前3位、后1位不加密
	if len(t.StdAreaCode())+n < TelMinLength || num[0] == '+' {
		return "", stdfmt.TelZero, fmt.Errorf("tel %s is too short or contains nation code", tel)
	}

	cipherStart := 3   // 前三位不加密，兼顾加密和维护效率
	cipherEnd := n - 1 // 最后1位不加密
	var seedChar byte
	var s strings.Builder
	s.Grow(len(tel)) // 预估加密后长度
	ac := t.StdAreaCode()
	if nc := t.StdNationCode(); nc != "" {
		// 省略中国国际区号
		if nc != "86" {
			nc, err = types.ConvertBase(nc, base10, base64) // 用64进制数字编码，并用分隔符隔开
			if err != nil {
				return "", stdfmt.TelZero, fmt.Errorf("invalid std nation code %s", nc)
			}
			s.WriteString(nc)
		}
		s.WriteByte(TelCipherSeparator)
		if t.StdAreaCode() == "" {
			s.WriteByte(TelCipherSeparator)
		}
	}
	if ac != "" {
		// 区号base36应该是 1-2个字符
		ac, err = types.ConvertBase(t.StdAreaCode(), base10, base64) // 用64进制数字编码，并用分隔符隔开
		if err != nil {
			return "", stdfmt.TelZero, fmt.Errorf("invalid std area code %s", ac)
		}
		s.WriteString(ac) // 区号不加密，使用36进制编码
		s.WriteByte(TelCipherSeparator)
		cipherStart = 0
	} else {
		ac, err := types.ConvertBase(num[:cipherStart], base10, base64) // 手机号前3位用64进制编码
		if err != nil || len(ac) != 2 {
			return "", stdfmt.TelZero, fmt.Errorf("invalid std area code %s", t.StdAreaCode())
		}
		s.WriteString(ac)
	}
	seed := int(num[0])
	keyClone := bytes.Clone(key)
	// 混淆key，会导致每次加密的结果都不一样，因此不适用于手机号唯一主键
	// 混淆key一般用于写之后，只查看的状况，因此可以多加一位
	if scatter {
		seed = rand.IntN(len(keyClone))
		seedChar = keyClone[seed]

		// 把key重组
		keyClone = coding.Scatter(keyClone, seed)
	}
	shift, ciphertext, err := encodeTelBase64Digits(num[cipherStart:cipherEnd], keyClone, seed, scatter)
	if err != nil {
		return "", stdfmt.TelZero, err
	}
	if err = coding.ShuffleEncrypt(ciphertext, shift, keyClone); err != nil {
		return "", stdfmt.TelZero, err
	}
	// 写入shiftChar、scatterShiftChar、密文、最后1位明文
	s.WriteByte(key[shift]) // shiftChar 是从原key中获取的
	if scatter {
		s.WriteByte(seedChar) // 如果shift落在偶数位，则第二位为 seedChar
	}
	s.Write(ciphertext)
	s.WriteByte(num[cipherEnd])
	str := s.String()
	minLen := TelCipherMinLen
	maxLen := TelCipherMaxLen
	if scatter {
		minLen = TelCipherScatterMinLen
		maxLen = TelCipherScatterMaxLen
	}
	if len(str) < minLen || len(str) > maxLen {
		log.Printf("[warn] ShuffleEncryptTel(%s, scatter:%v) length %d is out of [%d, %d]", tel, scatter, len(str), minLen, maxLen)
	}
	return TelCipher(str), t, nil

}
func encodeTelBase64Digits(text string, key []byte, nShift int, scatter bool) (int, []byte, error) {
	var (
		zeroPrefixNum int // 前缀0个数
		shift         int
		shiftMin      int
		shiftMax      int
	)
	for _, char := range text {
		if char != '0' {
			break
		}
		zeroPrefixNum++
	}
	var segValue uint64
	var err error
	s := text[zeroPrefixNum:]
	if s != "" {
		segValue, err = types.ParseUint64(s) // 去掉前缀0
		if err != nil {
			return 0, nil, ae.ErrInvalidInput
		}
	}

	// 填充了0 ，就会很短。因此，增加一位标识多少个0，很重要
	kn := len(key) / telShiftSegmentLen // 每节数字数量
	var shortZero bool
	if kn < telShiftSegmentItemMinLen {
		half := len(key) / 2
		rest := len(key) - half
		// 采用1/2

		if zeroPrefixNum > 0 {
			shiftMin = 0 // 前缀为0，放到前1/2
			shiftMax = rest
			shortZero = true
		} else {
			shiftMin = rest
			shiftMax = len(key) - 1
		}
		shift = int(segValue+uint64(nShift))%half + shiftMin

	} else {
		// 前缀0的个数，融入到shift 10个分段内。
		shiftMin = zeroPrefixNum * kn
		shiftMax = shiftMin + kn - 1
		shift = int((segValue+uint64(nShift))%uint64(kn)) + shiftMin
	}

	modWant := 1 // 奇数
	if scatter {
		modWant = 0 // 偶数
	}
	for shift%2 != modWant {
		shift++
		if shift >= shiftMax {
			shift = shiftMin
		}
	}

	// 连续3个0，则要把第3个0变为1（scatter=true时，用随机），再转为36进制。否则太短，导致容易侦破
	if zeroPrefixNum > 2 {
		varPrefix := byte('1') // 没有意义的字符，可以随意设置
		if scatter {
			varPrefix = byte(rand.UintN(9)) + '1' // 0-8 + '1' ==> '1'-'9'
		}
		variant := string(varPrefix) + text[3:] // 000xxxx ==> 001xxx
		if segValue, err = types.ParseUint64(variant); err != nil {
			return 0, nil, ae.ErrInvalidInput
		}
	}

	m := types.FormatBase64Uint(segValue)
	if shortZero {
		index := zeroPrefixNum + shift
		if index > len(key) {
			index -= len(key) // 相当于取余，但是减法性能比取余更好
		}
		// 前缀0的个数，放到第一位
		m = string(key[index]) + m
	}

	return shift, []byte(m), nil
}

func (t TelCipher) split(cipherTel TelCipher) (nationCode, areaCode, prefix []byte, shiftChar byte, cipherPart []byte, last byte, err error) {
	var nationCodeHandled bool
	segs := bytes.Split([]byte(cipherTel), []byte{TelCipherSeparator})
	switch len(segs) {
	case 1:
		cipherPart = segs[0]
	case 2:
		areaCode = segs[0]
		cipherPart = segs[1]
	case 3:
		nationCode = segs[0]
		if len(nationCode) == 0 {
			nationCode = []byte{'8', '6'}
			nationCodeHandled = true
		}
		areaCode = segs[1]
		cipherPart = segs[2]
	default:
		err = errors.New("cipher tel format error")
	}
	if !nationCodeHandled && len(nationCode) > 0 {
		var nc string
		nc, err = types.ConvertBase(string(nationCode), base64, base10)
		if err != nil {
			return
		}
		nationCode = []byte(nc)
	}
	if len(cipherPart) < 3 {
		err = ae.ErrInputWrongLength
		return
	}
	last = cipherPart[len(cipherPart)-1]
	cipherPart = cipherPart[:len(cipherPart)-1]
	if len(areaCode) > 0 {
		var ac string
		ac, err = types.ConvertBase(string(areaCode), base64, base10)
		if err != nil {
			return
		}
		areaCode = []byte(ac)
	} else {
		var pc string
		pc, err = types.ConvertBase(string(cipherPart[:2]), base64, base10)
		if err != nil || len(pc) != 3 {
			err = fmt.Errorf("invalid cipher tel prefix: %s=>%s", string(cipherPart[:2]), pc)
			return
		}
		prefix = []byte(pc)
		cipherPart = cipherPart[2:]
	}
	shiftChar = cipherPart[0]
	cipherPart = cipherPart[1:]
	return
}

// Decrypt
func (t TelCipher) Decrypt(key []byte) (stdfmt.Tel, error) {
	ValidateKeyLen(key, TelKeyMinLen)
	if t == "" {
		return stdfmt.TelZero, nil
	}
	nationCode, areaCode, prefix, shiftChar, ciphertext, last, err := t.split(t)
	if err != nil {
		return stdfmt.TelZero, err
	}
	keyIndex := make(map[byte]int, len(key))
	for i, k := range key {
		keyIndex[k] = i
	}
	var s strings.Builder
	s.Grow(12)
	if len(prefix) > 0 {
		s.Write(prefix) // 手机号前三位是明文
	}

	shift, ok := keyIndex[shiftChar]
	if !ok {
		return stdfmt.TelZero, coding.ErrCipherKeyMissChar(shiftChar)
	}
	scatter := shift%2 == 0
	var seedChar byte
	if scatter {
		seedChar = ciphertext[0]
		ciphertext = ciphertext[1:]
		randShift, ok := keyIndex[seedChar]
		if !ok {
			return stdfmt.TelZero, coding.ErrCipherKeyMissChar(seedChar)
		}
		key = coding.Scatter(key, randShift) // 把key重组
	}

	if err = coding.ShuffleDecrypt(ciphertext, shift, key); err != nil {
		return stdfmt.TelZero, err
	}

	plaintext, err := decodeTelBase64Digits(ciphertext, keyIndex, shift)
	if err != nil {
		return stdfmt.TelZero, err
	}
	s.WriteString(plaintext)
	s.WriteByte(last)

	tel := stdfmt.NewTel(string(nationCode), string(areaCode), s.String(), false)
	return tel, nil

}
func decodeTelBase64Digits(s []byte, keyIndex map[byte]int, shift int) (string, error) {
	var (
		zeroPrefixNum int
		shortZero     bool
	)

	kn := len(keyIndex) / telShiftSegmentLen // 每节数字数量
	if kn < telShiftSegmentItemMinLen {
		// 前缀为0的，放到前1/2
		if shift/2 == 0 {
			index, ok := keyIndex[s[0]]
			if !ok {
				return "", fmt.Errorf("char %s is not in key", string(s[1]))
			}
			zeroPrefixNum = index - shift
			if zeroPrefixNum < 0 {
				zeroPrefixNum += len(keyIndex) // 相当于取余，性能更好
			}
			shortZero = true
		}
	} else {
		zeroPrefixNum = shift / kn
	}
	if shortZero {
		s = s[1:] // 第一位被0的个数占据
	}

	text, err := types.ConvertBase(string(s), base64, base10)
	if err != nil {
		return "", err
	}
	if zeroPrefixNum > 2 {
		text = "000" + text[1:] //去掉第一位无意义数字
	} else if zeroPrefixNum > 0 {
		text = strings.Repeat("0", zeroPrefixNum) + text
	}
	return text, nil
}
func (t TelCipher) String() string {
	return string(t)
}

// 将密文转为脱敏后的电话号码
func (t TelCipher) Desensitize() (string, error) {
	if t == "" {
		return "", nil
	}
	nationCode, areaCode, prefix, _, main, last, err := t.split(t)
	if err != nil {
		return "", err
	}
	var s strings.Builder
	s.Grow(18) // 预估长度   +886 789-12345678  可能多1个scatter
	if len(nationCode) > 0 {
		s.WriteByte('+')
		s.Write(nationCode)
		s.WriteByte(' ')
	}
	if len(areaCode) > 0 {
		if len(nationCode) == 0 {
			s.WriteByte('0')
		}
		s.Write(areaCode)
		s.WriteByte('-')
	}
	if len(prefix) > 0 {
		s.Write(prefix)
	}
	s.Write(bytes.Repeat([]byte{'*'}, len(main))) // 没有判断scatter，可能会多1个*，不过不影响
	s.WriteByte(last)
	return s.String(), nil
}
