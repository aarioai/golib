package sms

import (
	"context"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/aa/atype"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/sdk/sms/aliyun"
	"github.com/aarioai/golib/sdk/sms/mo"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type VericodeSMSRequest struct {
	PseudoId    string `json:"pseudo_id"`
	PeriodLimit int    `json:"period_limit"`

	Sid         uint64        `json:"sid"`       // 发送批次ID
	SignName    string        `json:"sign_name"` // 公司签名
	Country     aenum.Country `json:"country_code"`
	PhoneNumber string        `json:"phone_number"`
	TplId       string        `json:"tpl_id"`

	VericodeLen int `json:"vericode_len"`
}

func (s *Service) SendAndCacheAliyunVericode(ctx context.Context, r VericodeSMSRequest) *ae.Error {
	if r.PeriodLimit > 0 {
		_, ok := s.h.ApplySmsVericodeSendingPermission(ctx, r.Country, r.PhoneNumber, r.PeriodLimit)
		if !ok {
			return ae.New(ae.TooManyRequests)
		}
	}

	if r.VericodeLen <= 0 {
		r.VericodeLen = 4
	}

	// 短信发送存在延时性问题，每5分钟内，重复发送相同的验证码
	vericode, ok := s.h.LoadSmsVericode(ctx, r.Country, r.PhoneNumber, r.PseudoId)
	if !ok {
		// 10分钟内没有发送过验证码，重新生成
		//支付宝用的是4位验证码，发送验证码前置需要人机对抗，所以即使4位数验证码是万分之一概率碰到，但是人机对抗很耗时。
		// 增加一个验证码ID，可以增加暴力破解的难度
		vericode = coding.RandNum(r.VericodeLen)
	}

	if !s.h.CacheSmsVericode(ctx, r.Country, r.PhoneNumber, vericode, r.PseudoId) {
		return NewE("cache sms vericode failed")
	}
	req := aliyun.VericodeRequest{
		Sid:         r.Sid,
		SignName:    r.SignName,
		Country:     r.Country,
		PhoneNumber: r.PhoneNumber,
		TplId:       r.TplId,
		Vericode:    vericode,
	}
	return s.SendAliyunVericode(req)
}

func (s *Service) VerifySmsVericode(ctx context.Context, cn aenum.Country, phoneNum, vericode, pseudoId string) *ae.Error {
	t := mo.SmsVerifyLog{
		Id:       bson.ObjectID{},
		Country:  cn,
		PhoneNum: phoneNum,
		Vericode: vericode,
		Errmsg:   "",
		VerifyAt: atype.Now(s.loc),
	}

	vcode, ok := s.h.LoadAndDeleteVericode(ctx, cn, phoneNum, pseudoId)
	if !ok {
		t.Errmsg = "not found"
		smsVerifyLog <- t
		return ae.New(ae.PreconditionFailed, "vericode not found")
	}
	if vcode != vericode {
		t.Errmsg = "not match"
		smsVerifyLog <- t
		return ae.New(ae.PreconditionFailed, "vericode not match")
	}

	smsVerifyLog <- t
	return nil
}
