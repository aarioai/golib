package aliyun

import (
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/golib/sdk/sms/config"
)

type Aliyun struct {
	AccessKey    string
	AccessSecret string
	RegionId     string
}

func NewAliyun(accessKey, accessSecret string, regionId ...string) *Aliyun {
	region := afmt.First(regionId)
	if region == "" {
		region = config.AliyunDefaultRegionId
	}
	return &Aliyun{
		RegionId:     region,
		AccessKey:    accessKey,
		AccessSecret: accessSecret,
	}
}

type SmsRequest struct {
	Sid          uint64            `json:"sid"`       // 发送批次ID
	SignName     string            `json:"sign_name"` // 公司签名
	Country      aenum.Country     `json:"country_code"`
	PhoneNumbers []string          `json:"phone_numbers"` // 支持对多个手机号码发送短信，手机号码之间以半角逗号（,）分隔。上限为1000个手机号码。批量调用相对于单条调用及时性稍有延迟。
	TplId        string            `json:"tpl_id"`
	TplParams    map[string]string `json:"tpl_params"`
}

type VericodeRequest struct {
	Sid         uint64        `json:"sid"`       // 发送批次ID
	SignName    string        `json:"sign_name"` // 公司签名
	Country     aenum.Country `json:"country_code"`
	PhoneNumber string        `json:"phone_number"`
	TplId       string        `json:"tpl_id"`
	Vericode    string        `json:"vericode"`
}
