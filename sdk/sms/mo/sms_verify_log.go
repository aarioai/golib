package mo

import (
	"github.com/aarioai/airis-driver/driver/index"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/golib/sdk/sms/config"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// // 日志只能增加，不能修改/删除。如果存在先后顺序，需要分两张表
type SmsVerifyLog struct {
	Id        bson.ObjectID  `bson:"_id" json:"id"`
	Country   aenum.Country  `bson:"country" json:"country"`
	PhoneNum  string         `bson:"phone_num" json:"phone_num"`
	Vericode  string         `bson:"vericode" json:"vericode"`
	Errmsg    string         `bson:"errmsg" json:"errmsg"`
	VerifyAt  atype.Datetime `bson:"verify_at" json:"verify_at"`
	CreatedAt atype.Datetime `bson:"created_at" json:"created_at"`
}

func (t SmsVerifyLog) Table() string {
	return config.SmsLogDbPrefix + "sms_verify_log"
}

func (t SmsVerifyLog) Indexes() index.Indexes {
	return index.NewIndexes(
		index.Primary("_id"),
	)
}
