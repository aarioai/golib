package sms

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/golib/sdk/sms/aliyun"
	"github.com/aarioai/golib/sdk/sms/enum"
	"github.com/aarioai/golib/sdk/sms/mo"
)

// SendAliyunVericode
// Suggest: SendAndCacheAliyunVericode
func (s *Service) SendAliyunVericode(r aliyun.VericodeRequest) *ae.Error {
	var log mo.SmsVericodeLog
	if s.enableLog {
		log = mo.SmsVericodeLog{
			Sid:        r.Sid,
			Uid:        0,
			Broker:     enum.BrokerAliyun,
			Country:    r.Country,
			PhonesNums: atype.SepStrings(r.PhoneNumber),
			MsgTpl:     r.TplId,
			Vericode:   r.Vericode,
			AckBizid:   "",
			AckMsg:     "",
			SendStatus: enum.SendUnknown,
			SentAt:     atype.Now(s.loc),
			//CreatedAt:  atype.Now(s.loc),
		}
	}
	res, e := s.aliyun.SendVericode(r)
	if e != nil {
		log.SendStatus = enum.SendFailed
		log.AckBizid = res.RequestId
		log.AckMsg = "【" + r.SignName + "】" + e.Text()
	} else {
		log.AckBizid = res.BizId
	}
	if s.enableLog {
		vericodeLog <- log
	}
	return e
}
