package aliyun

import (
	"encoding/json"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"strconv"
	"strings"
)

func NewError(err error) *ae.Error {
	if err == nil {
		return nil
	}
	return ae.NewError(err)
}

func NewE(msg string, args ...any) *ae.Error {
	return ae.NewE(msg, args...)
}

// Send 阿里云 签名公司名 和 模板 可以混合使用
func (s *Aliyun) Send(r SmsRequest) (*dysmsapi.SendSmsResponse, *ae.Error) {
	if len(r.PhoneNumbers) == 0 {
		return nil, ae.NewE("miss phone to send aliyun sms, %s, %s", r.SignName, r.TplId)
	}
	var params []byte
	if len(r.TplParams) > 0 {
		params, _ = json.Marshal(r.TplParams)
	}
	client, err := dysmsapi.NewClientWithAccessKey(s.RegionId, s.AccessKey, s.AccessSecret)
	if err != nil {
		return nil, NewError(err)
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = strings.Join(r.PhoneNumbers, ",")
	request.SignName = r.SignName
	request.TemplateCode = r.TplId
	request.TemplateParam = string(params)
	request.OutId = strconv.FormatUint(r.Sid, 10)
	resp, err := client.SendSms(request)
	if err != nil {
		return nil, NewError(err)
	}
	if resp.Code != "OK" {
		return nil, NewE("%s (%s)", resp.Message, resp.Code)
	}
	return resp, nil
}

func (s *Aliyun) SendVericode(r VericodeRequest) (*dysmsapi.SendSmsResponse, *ae.Error) {
	key := r.VericodeName
	if key == "" {
		key = "vericode"
	}
	sr := SmsRequest{
		Sid:          r.Sid,
		SignName:     r.SignName,
		Country:      r.Country,
		PhoneNumbers: []string{r.PhoneNumber},
		TplId:        r.TplId,
		TplParams: map[string]string{
			key: r.Vericode,
		},
	}
	return s.Send(sr)
}
