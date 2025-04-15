package crypto

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/lib/code/stdfmt"
	"github.com/aarioai/golib/libenum"
	"strings"
)

type LawyerLicCipher SidCipher
type LawyerCertCipher SidCipher

const (
	LawyerLicCipherKeyMinLen  = 10 // 0-9
	LawyerCertCipherKeyMinLen = 10
)

// ValidateLawyerLicEncryptKeys 一般启动的时候，验证key是否有错误。启动的时候执行，因此不必在意性能消耗
func ValidateLawyerLicEncryptKeys[T string | []byte](keys ...T) {
	// 基于 SidCipher
	ValidateSidEncryptKeys(keys...)
}

func ValidateLawyerCertEncryptKeys[T string | []byte](keys ...T) {
	// 基于 SidCipher
	ValidateSidEncryptKeys(keys...)
}

// ShuffleEncryptLawyerLic 律师执业证（含律师工作证）号采用17位代码’具体排序说明如下：
// 　　第1位为执业证书文本种类代码’1代表律师执业证文本。
// 　　第2-3位为持证人执业机构所在的省（区、市）代码: 执行《中华人民共和国行政区划代码》（6Ｂ2260）·
// 　　第4-5位为持证人执业机构所在的市（地、州｀盟）或者直辖市的区（县）代码’执行《中华人民共和国行政区划代码》（0Ｂ2260）。
// 　　第6—9位为首次批准律师执业的年度代码·持证人终止执业后申请重新执业的’为批准重新执业的年度代码。
// 　　第10位为律师执业证类别代码（专职律师1、兼职律师、香港居民律师3、澳门居民律师4、台湾居民律师5、公职律师6、公司律师7、法律援助律师8、军队律师9）。
// 　　第11位为性别代码（男0’女1）。
// 　　第12ˉ17位为律师执业证序列号代码（为避免因律师流动、变更执业证类别等产生重号现象’序列号编制保证一名律师只有一个唯一的序列号’全国范围内互不重号’从许可执业到终止执业该序列号永远不变）。
// 　　以王宇律师为例’其律师执业证号为：11102200810000003
// 　　1（律师执业证文本）11（北京市）02（延庆县）2008（首次批准律师执业年度）1（专职律师）0（男）000003（假 设的序列号）． 因律师流动、变更执业证类别等需要更改律师执业证号的’只更改变化的要素，其他要素不变.
// 　　假设一,王宇律师从北京市延庆县转到广东省深训市执业’只将北京市（11）、延庆县（02）的代码分别更改为广东省（44）、深训市（03）的代码’其他不变。王宇转到深训市执业后的代码应当为：14403200810000003 ，1（律师执业证）44（广东省）03（深训市）2008（首次批准律师执业年度）1（专职律师）0（男）000003（假设的序列号）。
// 　　假设二,王宇律师由专职律师转为兼职律师’只将专职代码（1）更改为兼职代码（2）’其他不变°王宇律师由专职律师转为兼职律师后的代码应当为：11102200820000003 ，1（律师执业证）11（北京市）02（延庆县）2008（批准律师执业年度）2（兼职律师）0（男）000003（假设的序列号）。
func ShuffleEncryptLawyerLic(s string, key []byte) (libenum.LawyerLicType, atype.Dist, atype.Year, LawyerLicCipher, error) {
	ValidateKeyLen(key, LawyerLicCipherKeyMinLen)
	var err error
	if s, err = stdfmt.ValidateLawyerLic(s); err != nil {
		return 0, 0, 0, "", err
	}
	dist, err1 := types.ParseUint16(s[1:5])
	year, err2 := types.ParseUint16(s[5:9])
	licType, err3 := types.ParseUint8(s[9:11])
	if ae.FirstError(err1, err2, err3) != nil {
		return 0, 0, 0, "", ae.ErrInvalidInput
	}

	// 律师证往往需要公开查询使用，而且也不需要非常隐私，因此不需要scatter加密 --> 方便SQL查询
	// 序号可能从0开始，如0000009，就会很小，触发不加密。因此不用进行base转换，直接当成字符串处理
	ciphertext, err := ShuffleEncryptSid(s[11:], key, false)
	if err != nil {
		return 0, 0, 0, "", err
	}
	cryptogram := LawyerLicCipher(ciphertext)
	return libenum.LawyerLicType(licType), atype.Dist(dist), atype.Year(year), cryptogram, nil
}
func (c LawyerLicCipher) String() string {
	return string(c)
}
func (c LawyerLicCipher) Desensitize(t libenum.LawyerLicType, d atype.Dist, y atype.Year) string {
	if c == "" {
		return ""
	}
	text := SidCipher(c).Desensitize()
	return "1" + d.String() + y.String() + t.String() + text
}
func (c LawyerLicCipher) Decrypt(t libenum.LawyerLicType, d atype.Dist, y atype.Year, key []byte) (string, error) {
	if c == "" {
		return "", nil
	}
	code, err := SidCipher(c).Decrypt(key)
	if err != nil {
		return "", err
	}
	lic := "1" + d.String() + y.String() + t.String() + code
	return lic, nil
}

// ShuffleEncryptLawyerCert 法律资格证
// 第1位：[ABC]  三种证书类型
// 第2-5位：4位发放年份
// 第6-11位：6位地区
// 第12-15位：4位该地区当年证书顺序编号
func ShuffleEncryptLawyerCert(s string, key []byte) (libenum.LawyerCertType, atype.Year, atype.Distri, LawyerCertCipher, error) {
	var err error
	if s, err = stdfmt.ValidateLawyerCert(s); err != nil {
		return 0, 0, 0, "", err
	}
	certType, ok := libenum.NewLawyerCertType(s[0])
	if !ok {
		return 0, 0, 0, "", ae.ErrInvalidInput
	}
	year, err1 := atype.ParseYear(s[1:4])
	distri, err2 := atype.ParseDistri(s[4:10])
	if ae.FirstError(err1, err2) != nil {
		return 0, 0, 0, "", ae.ErrInvalidInput
	}
	// 方便SQL查询，不用scatter
	ciphertext, err := ShuffleEncryptSid(s[10:], key, false)
	if err != nil {
		return 0, 0, 0, "", err
	}
	return certType, year, distri, LawyerCertCipher(ciphertext), nil
}
func (c LawyerCertCipher) Decrypt(t libenum.LawyerCertType, y atype.Year, d atype.Distri, key []byte) (string, error) {
	if c == "" {
		return "", nil
	}
	text, err := SidCipher(c).Decrypt(key)
	if err != nil {
		return "", err
	}
	return t.String() + y.String() + d.String() + text, nil
}
func (c LawyerCertCipher) Desensitize(t libenum.LawyerCertType, y atype.Year, d atype.Distri) string {
	if c == "" {
		return ""
	}
	sl := ShuffleDecryptLen(len(c))
	ast := strings.Repeat("*", sl-1) //  这个很短，因此不用前缀
	return t.String() + y.String() + d.String() + ast + string(c[len(c)-1])
}
func (c LawyerCertCipher) String() string {
	return string(c)
}
