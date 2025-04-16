package config

import "time"

var (
	SmsLogDbPrefix = "sdk_"

	AliyunDefaultRegionId = "cn-hangzhou"

	CacheDelimiter              = ":"
	CacheVericodeKeyFormat      = "libsdk:sms:vericode:%s"
	CacheVericodeLimitKeyFormat = "libsdk:sms:vericode_limit:%s"
	VericodeTTL                 = 10 * time.Minute // 短信验证码时效
	//VericodeLimit               = 10               // 每个 VericodePeriodTTL 内，每个账号（相同设备下）最多能发送多少。< 0 表示不限制
	VericodePeriodTTL = 24 * time.Hour // 每个周期时间，如每24小时每个账号最多发10条
)
