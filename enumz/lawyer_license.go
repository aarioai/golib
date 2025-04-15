package enumz

import "github.com/aarioai/airis/pkg/types"

// 17 位
// 　第1位为执业证书文本种类代码’1代表律师执业证文本。
// 　　第2ˉ3位为持证人执业机构所在的省（区、市）代码: 执行《中华人民共和国行政区划代码》（6Ｂ2260）·
// 　　第4ˉ5位为持证人执业机构所在的市（地、州｀盟）或者直辖市的区（县）代码’执行《中华人民共和国行政区划代码》（0Ｂ2260）。
// 　　第6—9位为首次批准律师执业的年度代码·持证人终止执业后申请重新执业的’为批准重新执业的年度代码。
// 　　第10位为律师执业证类别代码（专职律师1、兼职律师、香港居民律师3、澳门居民律师4、台湾居民律师5、公职律师6、公司律师7、法律援助律师8、军队律师9）。
// 　　第11位为性别代码（男0’女1）。
// 　　第12ˉ17位为律师执业证序列号代码（为避免因律师流动、变更执业证类别等产生重号现象’序列号编制保证一名律师只有一个唯一的序列号’全国范围内互不重号’从许可执业到终止执业该序列号永远不变）。
// 　　以王宇律师为例’其律师执业证号为：11102200810000003
// 　　1（律师执业证文本）11（北京市）02（延庆县）2008（首次批准律师执业年度）1（专职律师）0（男）000003（假 设的序列号）． 因律师流动、变更执业证类别等需要更改律师执业证号的’只更改变化的要素，其他要素不变.
// 　　假设一,王宇律师从北京市延庆县转到广东省深训市执业’只将北京市（11）、延庆县（02）的代码分别更改为广东省（44）、深训市（03）的代码’其他不变。王宇转到深训市执业后的代码应当为：14403200810000003 ，1（律师执业证）44（广东省）03（深训市）2008（首次批准律师执业年度）1（专职律师）0（男）000003（假设的序列号）。
// 　　假设二,王宇律师由专职律师转为兼职律师’只将专职代码（1）更改为兼职代码（2）’其他不变°王宇律师由专职律师转为兼职律师后的代码应当为：11102200820000003 ，1（律师执业证）11（北京市）02（延庆县）2008（批准律师执业年度）2（兼职律师）0（男）000003（假设的序列号）。
type LawyerLicType uint8

const (
	LawyerLicGBAM      LawyerLicType = 0  // 大湾区律师
	LawyerLicGBAF      LawyerLicType = 1  // 大湾区律师
	LawyerLicProM      LawyerLicType = 10 // 专职律师 男。尾数为0，表示男；1 表示女
	LawyerLicProF      LawyerLicType = 11
	LawyerLicPartTimeM LawyerLicType = 20 // 兼职律师
	LawyerLicPartTimeF LawyerLicType = 21
	LawyerLicHongkongM LawyerLicType = 30 // 港
	LawyerLicHongkongF LawyerLicType = 31
	LawyerLicMacaoM    LawyerLicType = 40 // 澳
	LawyerLicMacaoF    LawyerLicType = 41
	LawyerLicTaiwanM   LawyerLicType = 50 // 台
	LawyerLicTaiwanF   LawyerLicType = 51
	LawyerLicServantM  LawyerLicType = 60 // 公职律师
	LawyerLicServantF  LawyerLicType = 61
	LawyerLicCompanyM  LawyerLicType = 70 //  公司律师 - 法务
	LawyerLicCompanyF  LawyerLicType = 71
	LawyerLicAidM      LawyerLicType = 80 // 法律援助律师
	LawyerLicAidF      LawyerLicType = 81
	LawyerLicMilitaryM LawyerLicType = 90 // 军队律师
	LawyerLicMilitaryF LawyerLicType = 91
)

var LawyerLicNames = map[LawyerLicType]string{
	LawyerLicGBAM:      "大湾区律师",
	LawyerLicGBAF:      "大湾区律师",
	LawyerLicProM:      "专职律师",
	LawyerLicProF:      "专职律师",
	LawyerLicPartTimeM: "兼职律师",
	LawyerLicPartTimeF: "兼职律师",
	LawyerLicHongkongM: "香港律师",
	LawyerLicHongkongF: "香港律师",
	LawyerLicMacaoM:    "澳门律师",
	LawyerLicMacaoF:    "澳门律师",
	LawyerLicTaiwanM:   "台湾律师",
	LawyerLicTaiwanF:   "台湾律师",
	LawyerLicServantM:  "公职律师",
	LawyerLicServantF:  "公职律师",
	LawyerLicCompanyM:  "公司律师",
	LawyerLicCompanyF:  "公司律师",
	LawyerLicAidM:      "法援律师",
	LawyerLicAidF:      "法援律师",
	LawyerLicMilitaryM: "军队律师",
	LawyerLicMilitaryF: "军队律师",
}

func (t LawyerLicType) Name() string {
	if name, ok := LawyerLicNames[t]; ok {
		return name
	}
	return "律师"
}
func (t LawyerLicType) String() string {
	switch t {
	case 0:
		return "00" // 大湾区律师
	case 1:
		return "01"
	default:
		return types.FormatUint8(uint8(t))
	}
}
