package auth

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/sdk/auth/adto"
	"github.com/aarioai/golib/sdk/auth/config"
)

// 这里最好用string,因为 []byte是指针
func packUserToken(secure bool, atoken string, rtoken string, expiresIn, tokenTTL int64, admin uint8, scope map[string]any, conflict bool) adto.Token {

	if admin > 0 {
		if scope == nil {
			scope = map[string]any{"admin": admin}
		} else {
			scope["admin"] = admin
		}
	}

	return adto.Token{
		AccessToken: atoken,
		Conflict:    conflict,

		ValidateAPI: config.ValidateAPI,
		// Bearer  --> 客户端上传header: Authorization: Bearer $access_token
		TokenType:    config.UserTokenType,
		ExpiresIn:    expiresIn,
		RefreshToken: rtoken,
		RefreshAPI:   config.RefreshAPI,
		RefreshTTL:   tokenTTL,
		Secure:       secure,

		Scope: scope,
	}
}

// 一定要用string，用 []byte 是指针
func (s *Service) encryptUserFreshToken(atoken string) (string, *ae.Error) {
	// 如果用指针需要深度复制，而string是const，强转为 []byte 会重新开辟新的内存空间
	rtoken := []byte(atoken)

	for i, b := range rtoken {
		// 下划线是 95，小写字符是从97开始；A-Z 是 65-90；  数字0-9是：48-57
		if b < 91 && b > 64 {
			rtoken[i] = '_'
		}
	}
	base := []byte(config.UserTokenShuffleBase)
	if err := coding.ShuffleEncrypt(rtoken, userTokenBaseShift, base); err != nil {
		return "", ae.NewError(err)
	}
	for i, by := range rtoken {
		if by == '_' {
			rtoken[i] = 'A' + byte(i)%26 // 'A'=65, md5 是 base32， A-Z 是 26个字母
		}
	}
	return string(rtoken), nil
}

// 一定要用string，用 []byte 是指针
func (s *Service) DecryptUserFreshToken(rtoken string) (string, *ae.Error) {
	// 如果用指针需要深度复制，而string是const，强转为 []byte 会重新开辟新的内存空间
	atoken := []byte(rtoken)
	for i, b := range atoken {
		// 下划线是 95，小写字符是从97开始；A-Z 是 65-90；  数字0-9是：48-57
		if b < 91 && b > 64 {
			atoken[i] = '_'
		}
	}
	base := []byte(config.UserTokenShuffleBase)
	if err := coding.ShuffleDecrypt(atoken, userTokenBaseShift, base); err != nil {
		return "", ae.NewError(err)
	}
	atoken = replaceUserTokenUnderlines(atoken)
	return string(atoken), nil
}
