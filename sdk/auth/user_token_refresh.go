package auth

import (
	"context"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/aa/httpsvr/request"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/lib/code/coding"
	"github.com/aarioai/golib/sdk/auth/configz"
	"github.com/aarioai/golib/sdk/auth/dtoz"
	"github.com/aarioai/golib/typez"
	"time"
)

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
	base := []byte(configz.UserTokenShuffleBase)
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

// DecryptUserFreshToken
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
	base := []byte(configz.UserTokenShuffleBase)
	if err := coding.ShuffleDecrypt(atoken, userTokenBaseShift, base); err != nil {
		return "", ae.NewError(err)
	}
	atoken = replaceUserTokenUnderlines(atoken)
	return string(atoken), nil
}

// RefreshUserToken
// 微信access token 规则：
//  1. 若access_token已超时，那么进行refresh_token会获取一个新的access_token，新的超时时间；
//  2. 若access_token未超时，那么进行refresh_token不会改变access_token，但超时时间会刷新，相当于续期access_token。
//     refresh_token拥有较长的有效期（30天），当refresh_token失效的后，需要用户重新授权，所以，请开发者在refresh_token即将过期时（如第29天时），进行定时的自动刷新并保存好它。
//
// 这里规则：
//  1. access token 超时时间大于一半时间，那么返回之前的access token（时间也是之前的）
//  2. access token 超时时间小于一半时间（或超时），会拿到一个新的access token，并且刷新 user info 缓存
func (s *Service) RefreshUserToken(ctx context.Context, rtoken, currentPsid string, checkUser func(ctx context.Context, uid uint64) (typez.AdminLevel, *ae.Error)) (*dtoz.Token, *ae.Error) {
	atoken, e := s.DecryptUserFreshToken(rtoken)
	if e != nil {
		return nil, e
	}
	//  token 里面的 auth_at  是这个token生成时间，	SimplerUser 里面的 auth_at<返回客户端生成的那个 auth_at 就是这个> 是用户通过账号登陆时候的， 不是一回事
	svc, uid, vuid, ua, psid, authAt, factor, secureLogin, e := s.decryptUserToken(ctx, atoken)
	if e != nil {
		if e.Code == ae.Unauthorized {
			e = ae.ErrorPageExpired
		}
		return nil, e
	}
	if currentPsid != psid {
		return nil, ae.ErrorForbidden
	}
	cachedFactor, _ := s.h.LoadUserTokenFactor(ctx, svc, uid, ua)
	if factor != cachedFactor {
		return nil, ae.ErrorPageExpired
	}

	// refresh token 过期，则无效了
	refreshTokenExp := authAt + configz.UserRefreshTokenTTLs
	now := time.Now().Unix()
	if now > refreshTokenExp {
		return nil, ae.ErrorPageExpired
	}
	
	admin, e := checkUser(ctx, uid)
	if e != nil {
		return nil, e
	}

	expiresIn := authAt + configz.UserRefreshTokenTTLs - now // token 剩余时间（秒）
	if expiresIn > configz.UserTokenWillRefresh {
		rtokenExpiresIn := refreshTokenExp - now
		t := PackUserToken(secureLogin, atoken, rtoken, expiresIn, rtokenExpiresIn, admin, nil, false)
		return &t, nil
	}

	return s.NewUserToken(ctx, svc, uid, vuid, ua, psid, admin, false, secureLogin)
}

// GrantUserToken
// 客户端判断token临过期时间，主动往服务器通过refresh token换取新的 access token
// 使用GET
// https://developers.weixin.qq.com/doc/oplatform/Website_App/WeChat_Login/Authorized_Interface_Calling_UnionID.html
func (s *Service) GrantUserToken(ctx context.Context, r *request.Request, psid string, checkUser func(ctx context.Context, uid uint64) (typez.AdminLevel, *ae.Error)) (token *dtoz.Token, e *ae.Error) {
	// https://www.oauth.com/oauth2-servers/access-tokens/refreshing-access-tokens/
	// https://community.atlassian.com/t5/Bitbucket-questions/How-we-can-get-a-refresh-token-with-rest-api/qaq-p/793325
	// curl -X POST -u "{client_id:secret}" https://bitbucket.org/site/oauth2/access_token -d grant_type=authorization_code -d code={code}
	// curl -X POST -u "{client_id}:{secret}" https://bitbucket.org/site/oauth2/access_token -d grant_type=refresh_token -d refresh_token={refresh_token}
	// authorization_code（拿code换） or refresh_token
	// 使用 refresh_token，客户端就需要传递 client_id 和 client_secret；用 authorization_code ，就不需要传递
	// client_id = uid,  client_secret = access token  不建议！
	grantType, e0 := r.Body(enumz.ParamGrantType, `^(authorization_code|refresh_token)$`)
	//code, e1 := r.Body("code", false)                  // 用于authorization_code（拿code换）授权
	code, e1 := r.Body(enumz.ParamRefreshToken) // 用于更新
	// scope, e2 := r.Body(conf.ParamScope, false )
	if e = ae.First(e0, e1); e != nil {
		return nil, e
	}
	switch grantType.String() {
	case "authorization_code":
	case "refresh_token":
		return s.RefreshUserToken(ctx, code.String(), psid, checkUser)
	}
	return nil, ae.NewBadParam(enumz.ParamGrantType)
}
