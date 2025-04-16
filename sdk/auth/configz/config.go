package configz

var (
	UserTokenCryptMd5Key = ""
	UserTokenShuffleBase = ""

	UserTokenType = "Bearer"
	ValidateAPI   = "HEAD /api/v1/pas/auth/access_token"
	RefreshAPI    = "PUT /api/v1/pas/auth/access_token"

	// 不要使用MySQL存储 RAS key，及一切带换行符的key
	UserTokenCryptMd5key    = "libsdk_auth.cipher_user_token_md5key"
	UserTokenShuffleBaseKey = "libsdk_auth.cipher_user_token_shuffle_base"
	MmcCryptSecret          = "a_pas.cipher.user_mmc_key"

	// rsa 文件 包括以下几种
	// 1.  b/s 架构     JS 可用的    bs768  bs1024 bs2048  bs4096
	// 2.  c/s 架构     客户端可用
	// 3.  s/s 架构     服务端用
	// xxx.rsa 是密钥；  xxx.rsa.pub 是公开密钥
	UserPasswordRSAPubkeyDERB64 = "rsa.app-2048.pub.der.b64"
	UserPasswordRSAPrivkeyDER   = "rsa.app-2048.priv.der"

	MmcRSAPubkeyDERB64 = "rsa.app-512.pub.der.b64"
	MmcRSAPrivkeyDER   = "rsa.app-512.priv.der"

	CachePrefix = "libsdk:auth:"

	UserTokenTimeWindow   = int64(10 * 60)         // 10分钟的时间窗口
	UserTokenTTLs         = int64(12 * 3600)       //  time-to-life in seconds 客户端token 只保留12小时
	UserRefreshTokenTTLs  = int64(181 * 24 * 3600) //  time-to-life in seconds refresh token 有效期
	UserTokenIntervalDays = uint8(7)
	UserTokenWillRefresh  = 4 * 3600 // 若access token剩余小于这个值，使用refresh token才会更新；否则会沿用之前的
)
