package crypto

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"github.com/aarioai/golib/libenum"
)

type UsccCipher SidCipher

// ValidateUSCCEncryptKeys 一般启动的时候，验证key是否有错误。启动的时候执行，因此不必在意性能消耗
func ValidateUSCCEncryptKeys[T string | []byte](keys ...T) {
	// 基于 SidCipher
	ValidateSidEncryptKeys(keys...)
}

// ShuffleEncryptUSCC
// 第1位：[\dA-Z]登记管理部门代码，表示该组织由哪个部门进行登记管理。
// 第2位：机构类别代码，表示该组织的机构类型。
// 第3-8位：登记管理机关行政区划码，表示该组织登记管理机关所在的行政区划。
// 第9-17位：主体标识码，用于唯一标识一个法人或其他组织。
// 第18位：校验码，用于校验整个代码的正确性。
func ShuffleEncryptUSCC(s string, key []byte) (t libenum.UsccType, distri atype.Distri, cryptogram UsccCipher, err error) {

	if s, err = stdfmt.ValidateUSCC(s); err != nil {
		return
	}
	var ok bool
	t, ok = libenum.ToUsccType(s[0:2])
	if !ok {
		err = ae.ErrInvalidInput
		return
	}

	distri, err = atype.ParseDistri(s[2:8])
	if err != nil {
		err = ae.ErrInvalidInput
		return
	}
	var ciphertext SidCipher
	ciphertext, err = ShuffleEncryptSid(s[8:], key, false) // 需要查询，不应该scatter
	if err != nil {
		return
	}
	cryptogram = UsccCipher(ciphertext)
	return
}

func (c UsccCipher) Desensitize(t libenum.UsccType, distri atype.Distri) string {
	if c == "" {
		return ""
	}
	text := SidCipher(c).Desensitize()
	return t.Code() + distri.String() + text
}

func (c UsccCipher) String() string {
	return string(c)
}

func (c UsccCipher) Decrypt(t libenum.UsccType, distri atype.Distri, key []byte) (string, error) {
	if c == "" {
		return "", nil
	}
	text, err := SidCipher(c).Decrypt(key)
	if err != nil {
		return "", err
	}
	return t.Code() + distri.String() + text, nil
}
