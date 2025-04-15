package crypto

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"strings"

	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/aa/atype"
)

type CinCipher string

const (
	CinCipherKeyMinLen = 11 // [0-9X]
)

// ValidateCinEncryptKeys 一般启动的时候，验证key是否有错误。启动的时候执行，因此不必在意性能消耗
func ValidateCinEncryptKeys[T string | []byte](keys ...T) {
	// 基于FastEncrypt
	coding.ValidateShuffleEncryptKeys(keys...)
}

// ShuffleEncryptCin
// distri 保持不变，方便统计
// birth date  基本保持不变，方便统计
// sex 方便统计
// （1）前1、2位数字表示：所在省份的代码；
// （2）第3、4位数字表示：所在城市的代码；
// （3）第5、6位数字表示：所在区县市（县级市）的代码；
// （4）第7~14位数字表示：出生年、月、日；7.8.9.10位是年，11.12位是月13.14位是日
// （5）第15、16、17位是顺序码，顺序码是表示同一地址码所标识的区域范围内，对同年、同月、同日出生的人编定的顺序号。第17位数字表示性别：奇数表示男性，偶数表示女性；
// （6）第18位数字是校检码：校检码可以是0~9的数字，有时也用X表示。
func ShuffleEncryptCin(cin string, key []byte) (distri atype.Distri, birthDate atype.Date, sex aenum.Sex, cryptogram CinCipher, err error) {
	ValidateKeyLen(key, CinCipherKeyMinLen)

	if cin, err = stdfmt.ValidateCIN(cin); err != nil {
		return
	}

	distri, err = atype.ParseDistri(cin[0:6])
	_, err2 := types.ParseUint(cin[6:14])
	if err = ae.FirstError(err, err2); err != nil {
		return
	}
	birthDate = atype.Date(cin[6:10] + "-" + cin[10:12] + "-" + cin[12:14])
	sex = aenum.Male
	if (cin[16]-'0')%2 == 0 {
		sex = aenum.Female
	}

	lastN := len(cin) - 1
	last := CinCipher(cin[lastN])
	id := cin[14:lastN]
	var ciphertext TextCipher
	// 身份证号以后可能需要用到查询。scatter设为false。非常短，因此直接当字符串处理
	if ciphertext, err = ShuffleEncrypt(id, key, false); err != nil {
		return
	}
	cryptogram = CinCipher(ciphertext) + last
	return
}
func (c CinCipher) String() string {
	return string(c)
}
func (c CinCipher) Desensitize(distri atype.Distri, birthDate atype.Date) string {
	if c == "" {
		return ""
	}
	sl := ShuffleDecryptLen(len(c))
	lastN := len(c) - 1
	last := string(c[lastN])
	b := strings.ReplaceAll(birthDate.String(), "-", "")
	ast := strings.Repeat("*", sl-1)
	return distri.String() + b + ast + last
}
func (c CinCipher) Decrypt(distri atype.Distri, birthDate atype.Date, key []byte) (string, error) {
	if len(c) < ShuffleEncryptMinLen+1 {
		return "", ae.ErrInvalidInput
	}
	lastN := len(c) - 1
	last := string(c[lastN])
	ciphertext := []byte(c[:lastN])
	plaintext, err := TextCipher(ciphertext).Decrypt(key)
	if err != nil {
		return "", err
	}
	b := strings.ReplaceAll(birthDate.String(), "-", "")
	return distri.String() + b + plaintext + last, nil
}
