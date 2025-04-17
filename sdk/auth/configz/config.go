package configz

var (
	UserTokenCryptMd5Key = ""
	UserTokenShuffleBase = ""

	UserTokenType = "Bearer"
	ValidateAPI   = "HEAD /api/v1/pas/auth/access_token"
	RefreshAPI    = "PUT /api/v1/pas/auth/access_token"
	
	CachePrefix = "libsdk:auth:"

	UserTokenTimeWindow   = int64(10 * 60)         // 10分钟的时间窗口
	UserTokenTTLs         = int64(12 * 3600)       //  time-to-life in seconds 客户端token 只保留12小时
	UserRefreshTokenTTLs  = int64(181 * 24 * 3600) //  time-to-life in seconds refresh token 有效期
	UserTokenIntervalDays = uint8(7)
	UserTokenWillRefresh  = 4 * 3600 // 若access token剩余小于这个值，使用refresh token才会更新；否则会沿用之前的
)
