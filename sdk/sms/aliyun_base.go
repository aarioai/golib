package sms

import (
	"context"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/golib/sdk/sms/aliyun"
	"github.com/aarioai/golib/sdk/sms/enum"
	"github.com/aarioai/golib/sdk/sms/mo"
)

// SendAliyunVericode
// Suggest: SendAndCacheAliyunVericode
func (s *Service) SendAliyunVericode(ctx context.Context, r aliyun.VericodeRequest) *ae.Error {
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
			RequestId:  "",
			BizId:      "",
			AckMsg:     "",
			SendStatus: enum.SendUnknown,
			SendAt:     atype.Now(s.loc),
			//CreatedAt:  atype.Now(s.loc),
		}
	}
	res, e := s.aliyun.SendVericode(r)
	log.RequestId = res.RequestId
	log.BizId = res.BizId
	if !s.app.Check(ctx, e) {
		log.SendStatus = enum.SendFailed
		log.AckMsg = "【" + r.SignName + "】" + e.Text()
	} else {
		log.SendStatus = enum.SendOK
	}
	if s.enableLog {
		smsVericodeLog <- log
	}
	return e
}
