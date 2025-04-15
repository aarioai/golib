package stdfmt

import (
	"github.com/aarioai/airis/aa/ae"
	"strings"
)

const (
	CinLen           = 18
	UsccLen          = 18
	LawyerLicMinLen  = 17 // 律师序号可能要满了，暂时保持兼容18位
	LawyerCertMinLen = 15
)

// ValidateCIN 中国身份证号码
func ValidateCIN(s string) (string, error) {
	s = strings.ToUpper(ReplaceWithStdASCII(s, false))
	if len(s) != CinLen {
		return "", ae.ErrInputWrongLength
	}
	lastN := len(s) - 1
	last := s[lastN]
	for i := 0; i < lastN-1; i++ {
		if s[i] < '0' || s[i] > '9' {
			return "", ae.ErrInvalidInput
		}
	}
	if last != 'X' && (last < '0' || last > '9') {
		return "", ae.ErrInvalidInput
	}
	return s, nil
}

// ValidateUSCC 统一社会信用代码
// 第1位：[\dA-Z]登记管理部门代码，表示该组织由哪个部门进行登记管理。
// 第2位：机构类别代码，表示该组织的机构类型。
// 第3-8位：登记管理机关行政区划码，表示该组织登记管理机关所在的行政区划。
// 第9-17位：主体标识码，用于唯一标识一个法人或其他组织。
// 第18位：校验码，用于校验整个代码的正确性。
func ValidateUSCC(s string) (string, error) {
	s = strings.ToUpper(ReplaceWithStdASCII(s, false))
	if len(s) != UsccLen {
		return "", ae.ErrInputWrongLength
	}
	for _, c := range s {
		if c < '0' || (c > '9' && c < 'A') || c > 'Z' {
			return "", ae.ErrInvalidInput
		}
	}
	return s, nil
}

// ValidateLawyerLic 律师执业证（含律师工作证）号采用17位代码’具体排序说明如下：
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
func ValidateLawyerLic(s string) (string, error) {
	s = ReplaceWithStdNumbers(s)
	if len(s) < LawyerLicMinLen {
		return "", ae.ErrInputWrongLength
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return "", ae.ErrInvalidInput
		}
	}
	return s, nil
}

// ValidateLawyerCert 法律资格证
// 第1位：[ABC]  三种证书类型
// 第2-5位：4位发放年份
// 第6-11位：6位地区
// 第12-15位：4位该地区当年证书顺序编号
func ValidateLawyerCert(s string) (string, error) {
	s = strings.ToUpper(ReplaceWithStdASCII(s, false))
	if len(s) < LawyerCertMinLen {
		return "", ae.ErrInputWrongLength
	}
	if s[0] < 'A' || s[0] > 'C' {
		return "", ae.ErrInvalidInput
	}
	for i := 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return "", ae.ErrInvalidInput
		}
	}
	return s, nil
}
