package crypto

// 身份证号脱敏
func DesCin(idNum string) string {
	if len(idNum) == 18 {
		return idNum[0:3] + "*************" + idNum[16:]
	}
	return idNum
}

// 统一信用代码脱敏
//
//	保留前10位、后4位
func DesUscc(idNum string) string {
	if len(idNum) == 18 {
		return idNum[0:10] + "****" + idNum[14:]
	}
	return idNum
}

// 姓名脱敏
func DesChineseName(name string) string {
	n := []rune(name)
	if len(n) == 0 {
		return name
	}
	return string(n[0]) + "*" + string(n[len(n)-1:])
}
func DesWeixin(weixin string) string {
	if weixin == "" {
		return ""
	}
	l := len(weixin)
	if l < 2 {
		return "***"
	}
	return weixin[0:1] + "***" + weixin[l-1:]
}

// 法律职业资格证号脱敏
func DesLawyerLicense(num string) string {
	l := len(num)
	if l < 11 {
		return "********"
	}
	return num[0:9] + "******" + num[l-3:]
}

// 执业证号脱敏
func DesLawyerCert(num string) string {
	l := len(num)
	if l < 8 {
		return "********"
	}
	return num[0:5] + "*******" + num[l-3:]
}
