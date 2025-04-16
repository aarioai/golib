package mo

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/golib/sdk/sms/configz"
	"github.com/aarioai/golib/sdk/sms/enum"
)

// 短信验证码日志单独记录吧
// 日志只能增加，不能修改/删除。如果存在先后顺序，需要分两张表
type SmsVericodeLog struct {
	Sid        uint64           `bson:"sid" json:"sid"` // sms id
	Uid        uint64           `bson:"uid" json:"uid"`
	Broker     enum.SmsBroker   `bson:"broker" json:"broker"`
	Country    aenum.Country    `bson:"country" json:"country"`
	PhonesNums atype.SepStrings `bson:"phone_nums" json:"phone_nums"`
	MsgTpl     string           `bson:"msg_tpl" json:"msg_tpl"`
	Vericode   string           `bson:"vericode" json:"vericode"`
	AckBizid   string           `bson:"ack_bizid" json:"ack_bizid"`
	AckMsg     string           `bson:"ack_msg" json:"ack_msg"`
	SendStatus enum.SendStatus  `bson:"send_status" json:"send_status"`
	SendAt     atype.Datetime   `bson:"send_at" json:"send_at"`
	CreatedAt  atype.Datetime   `bson:"created_at" json:"created_at"`
}

func (t SmsVericodeLog) Table() string {
	return configz.SmsLogDbPrefix + "sms_vericode_log"
}

func (t SmsVericodeLog) Indexes() index.Indexes {
	return index.NewIndexes(
		index.Primary("sid"),
	)
}
