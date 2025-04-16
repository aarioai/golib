package configz

import "time"

var (
	SmsLogDbPrefix = "sdk_"

	AliyunDefaultRegionId = "cn-hangzhou"

	CachePrefix = "libsdk:sms:"
	VericodeTTL = 10 * time.Minute // 短信验证码时效
	//VericodeLimit               = 10               // 每个 VericodePeriodTTL 内，每个账号（相同设备下）最多能发送多少。< 0 表示不限制
	VericodePeriodTTL = 24 * time.Hour // 每个周期时间，如每24小时每个账号最多发10条
)
