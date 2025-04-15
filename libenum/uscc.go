package libenum

import (
	"encoding/json"
	"github.com/aarioai/airis/pkg/types"
)

type UsccType uint8
type UsccStatus int8

const (
	UsccDelete        UsccStatus = -128 // 注销（吊销、停业、清算）
	UsccStatusUnknown UsccStatus = 0    //其他
	UsccInBusiness    UsccStatus = 1    // 续存（在营、开业、在册）

)

const (
	UsccInstOrg       UsccType = 11  //  机关
	UsccInstUnit      UsccType = 12  // 事业单位
	UsccInst          UsccType = 13  // 其他机构编制
	UsccFaNews        UsccType = 21  // 外国常驻新闻机构
	UsccFa            UsccType = 22  // 其他外交机构
	UsccLawfirm       UsccType = 31  // 律师执业机构
	UsccAdminNotary   UsccType = 32  // 公证处
	UsccLegalSvc      UsccType = 33  // 基层法律服务所
	UsccJudi          UsccType = 34  // 司法鉴定机构
	UsccArbi          UsccType = 35  // 仲裁委员会
	UsccAdmin         UsccType = 39  // 其他司法行政
	UsccCulForeign    UsccType = 41  // 外国在华文化中心
	UsccCul           UsccType = 49  // 其他文化
	UsccSocialGroup   UsccType = 51  // 社会团体
	UsccPrivateUnit   UsccType = 52  // 民办非企业单位
	UsccFoundation    UsccType = 53  // 基金会
	UsccCivil         UsccType = 59  // 其他民政
	UsccTourForeign   UsccType = 61  // 外国旅游部门常驻代表机构
	UsccTourGAT       UsccType = 62  // 港澳台地区旅游部门常驻内地（大陆）代表机构
	UsccTour          UsccType = 69  // 其他旅游
	UsccWorshipPlace  UsccType = 71  // 宗教活动场所
	UsccWorshipSchool UsccType = 72  //宗教院校
	UsccWorship       UsccType = 79  // 其他宗教
	UsccUnionBase     UsccType = 81  // 基层工会
	UsccUnion         UsccType = 89  // 其他工会
	UsccCompany       UsccType = 91  // 企业
	UsccIndividualBiz UsccType = 92  // 个体工商户
	UsccFarmCo        UsccType = 93  // 农民专业合作社
	UsccMilitaryUnit  UsccType = 101 // A1 军队事业单位
	UsccMilitary      UsccType = 109 // A9 其他 中央军委改革和编制办公室
	UsccCoEcoG        UsccType = 231 // N1 组级集体经济组织
	UsccCoEcoV        UsccType = 232 // N2 村级集体经济组织
	UsccCoEcoT        UsccType = 233 // N3 乡镇级集体经济组织
	UsccCoEco         UsccType = 239 // N9 其他农业集体经济组织
	UsccOther         UsccType = 251 // Y1 其他
)

func ToUsccType(s string) (UsccType, bool) {
	if len(s) != 2 {
		return 0, false
	}
	a := s[0]
	b, err := types.ParseUint8(s[1:])
	if err != nil {
		return 0, false
	}
	if a >= '1' && a <= '9' {
		n := (a-'0')*10 + b
		return UsccType(n), true
	} else if a >= 'A' && a <= 'O' {
		n := (a-'A'+10)*10 + b
		return UsccType(n), true
	} else if a == 'Y' {
		return UsccOther, true
	}
	return 0, true
}

func (t UsccType) Code() string {
	if t < 100 {
		return types.FormatUint8(uint8(t))
	}

	a := t / 10
	b := string(t%10 + '0')
	c := a - 10 + 'A'
	if c <= 'O' {
		return string(c) + b
	}
	return "Y1"
}

type CorpType uint8

const (
	InstOrg    CorpType = 1 // 机关
	Court      CorpType = 2 // 法院（未定- 包含最高院）
	InstUnit   CorpType = 3 // 事业单位
	Lawfirm    CorpType = 7 // 律所
	Enterprise CorpType = 9 // 企业
)

var (
	CorpTypeNames = map[CorpType]string{
		InstOrg:    "机关",
		Court:      "法院",
		InstUnit:   "事业单位",
		Lawfirm:    "律师事务所",
		Enterprise: "企业",
	}

	CorpTypesJson string
)

func init() {
	cs, _ := json.Marshal(CorpTypeNames)
	CorpTypesJson = string(cs)
}
