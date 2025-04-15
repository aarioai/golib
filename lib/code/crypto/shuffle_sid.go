package crypto

import (
	"github.com/aarioai/golib/lib/code/coding"
	"strings"
)

type SidCipher string // 使用于微信ID等 string id

var (
	SidEncryptMinLen = ShuffleEncryptMinLen + 2 // 首尾两个字符明文
)

// ValidateSidEncryptKeys 一般启动的时候，验证key是否有错误。启动的时候执行，因此不必在意性能消耗
func ValidateSidEncryptKeys[T string | []byte](keys ...T) {
	// 基于FastEncrypt
	coding.ValidateShuffleEncryptKeys(keys...)
}

// ShuffleEncryptSid 加密所有string id，首位和最后一位是明文
func ShuffleEncryptSid(s string, key []byte, scatter bool) (SidCipher, error) {
	if len(s) < SidEncryptMinLen {
		return SidCipher(s), nil // 字符太少，直接显示明文
	}

	first := SidCipher(s[0])
	lastN := len(s) - 1
	last := SidCipher(s[lastN])
	ciphertext, err := ShuffleEncrypt(s[1:lastN], key, scatter)
	if err != nil {
		return "", err
	}
	return first + SidCipher(ciphertext) + last, nil
}

func (c SidCipher) Decrypt(key []byte) (string, error) {
	if ShuffleDecryptLen(len(c)) < SidEncryptMinLen {
		return string(c), nil // 明文，直接显示
	}
	first := string(c[0])
	lastN := len(c) - 1
	last := string(c[lastN])
	ciphertext := TextCipher(c[1:lastN])
	plaintext, err := ciphertext.Decrypt(key)
	if err != nil {
		return "", err
	}
	return first + plaintext + last, nil
}
func (c SidCipher) Desensitize(wantLen ...int) string {
	if c == "" {
		return ""
	}
	sl := ShuffleDecryptLen(len(c))
	fixLen := sl
	if len(wantLen) > 0 && wantLen[0] > 2 {
		fixLen = wantLen[0]
	}
	if sl < 2 {
		return strings.Repeat("*", fixLen)
	}

	ast := strings.Repeat("*", fixLen-2)
	return string(c[0]) + ast + string(c[len(c)-1])
}

func (c SidCipher) String() string {
	return string(c)
}
