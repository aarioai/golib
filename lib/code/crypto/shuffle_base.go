package crypto

import (
	"bytes"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/lib/code/coding"
	"math/rand/v2"
	"strconv"
)

// shuffle 加密强度低，但是计算量少。使用需要频繁转换、且对加密要求不太高的场景

type DigitsCipher string
type TextCipher string

const (
	ShuffleEncryptMin    = 1000
	ShuffleEncryptMinLen = 2 // 1000=base64(fE)
)

func ValidateKeyLen(key []byte, minLen int) {
	if len(key) < minLen {
		panic("encrypt key is too short") // 可控的人为失误，应当直接panic
	}
}
func ShuffleDecryptLen(l int) int {
	if l == 0 {
		return 0
	}
	return l - 1 // 1位前缀shift
}

// ShuffleEncryptDigits 对数字shuffle加密
// key ：0-9  表示使用数字重新组合；
// key: 0-9a-f  表示使用16进制字符组合
// key :0-9a-z 表示数字+小写字母组合
// key : 64进制，表示64进制组合
func ShuffleEncryptDigits(n uint64, key []byte, scatter bool) (DigitsCipher, error) {
	ValidateKeyLen(key, base10)
	if n < ShuffleEncryptMin {
		return "", ae.ErrInputTooSmall // 数字太小，加密容易破解出key
	}
	var s string
	if len(key) > base36 {
		s = types.FormatBase64Uint(n)
	} else {
		s = strconv.FormatUint(n, len(key))
	}
	ciphertext, err := ShuffleEncrypt(s, key, scatter)
	if err != nil {
		return "", err
	}

	return DigitsCipher(ciphertext), nil
}

// ShuffleEncrypt  长度增加了1位，因此必须要重新开辟新内存，因此直接用 string,而不要过度优化使用 []byte
func ShuffleEncrypt(text string, key []byte, scatter bool) (TextCipher, error) {
	if len(text) < ShuffleEncryptMinLen {
		return "", ae.ErrInputTooShort // 数字太小，加密容易破解出key
	}
	var shift int
	var shiftChar byte
	modBase := len(key) - 2
	// 混淆key，会导致每次加密的结果都不一样，因此不适用于手机号唯一主键
	// 混淆key一般用于写之后，只查看的状况，因此可以多加一位
	if scatter {
		shift = rand.IntN(modBase)
		if shift%2 == 0 { // 保持奇数
			shift++
		}
		shiftChar = key[shift] // 必须要在使用没scatter之前，即原始key
		// 把key重组
		key = bytes.Clone(key)
		key = coding.Scatter(key, shift)
	} else {
		shift = int(text[0]) % modBase
		if shift%2 == 1 { // 保持偶数
			shift++
		}
		shiftChar = key[shift]
	}

	ciphertext := make([]byte, 1, len(text)+1)
	ciphertext[0] = shiftChar
	ciphertext = append(ciphertext, text...)
	if err := coding.ShuffleEncrypt(ciphertext[1:], shift, key); err != nil {
		return "", err
	}
	return TextCipher(ciphertext), nil
}

func (c DigitsCipher) Decrypt(key []byte) (uint64, error) {
	ValidateKeyLen(key, base10)
	if len(key) == 0 {
		return 0, nil
	}
	s, err := TextCipher(c).Decrypt(key)
	if err != nil {
		return 0, err
	}
	if len(key) > base36 {
		return types.ParseBase64Uint(s)
	}
	return strconv.ParseUint(s, len(key), 64)
}

func (c TextCipher) Decrypt(key []byte) (string, error) {
	if c == "" {
		return "", nil
	}
	if len(c) < ShuffleEncryptMinLen+1 {
		return "", ae.ErrInputTooShort // 数字太小，加密容易破解出key
	}
	shiftChar := c[0]
	shift := bytes.IndexByte(key, shiftChar)
	if shift < 0 {
		return "", coding.ErrCipherKeyMissChar(shiftChar)
	}
	scatter := shift%2 == 1
	if scatter {
		key = coding.Scatter(key, shift)
	}
	text := []byte(c[1:])
	if err := coding.ShuffleDecrypt(text, shift, key); err != nil {
		return "", err
	}
	return string(text), nil
}
