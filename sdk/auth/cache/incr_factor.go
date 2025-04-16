package cache

import (
	"context"
	"github.com/aarioai/golib/enumz"
	"github.com/aarioai/golib/sdk/auth/configz"
	"github.com/aarioai/golib/sdk/cachez"
	"github.com/aarioai/golib/typez"
	"strconv"
	"time"
)

// refresh token 跟 access token 数据基本一致，就是过期时间比access token 长一倍以上
// access token 不保存完整（节省内存），只保存自增ID加密因子。access token 有效期是2小时（客户端需要经常换refresh token），加密因子跟freshtoken ttl 一致，保持7天

// 加密算法本身可信，但是防止token重复使用、被盗用，还是需要加入一个cache因子

// svc uid  - 平台(ua)  唯一性    --> 保证登录的时候，不会把其他svc下的登录 factor值变了而导致退出
func toUserTokenFactorField(svc typez.Svc, uid uint64, ua enumz.UA) string {
	u := strconv.FormatUint(uid, 10)
	s := svc.String()
	if svc.Valid() {
		s = svc.String() + ":"
	}
	return s + "uid:" + u + ":ua:" + ua.String()
}
func incrUserTokenKeyPrefix() string {
	return configz.CachePrefix + "user_token_factor:"
}
func (h *Cache) IncrUserTokenFactor(ctx context.Context, svc typez.Svc, uid uint64, ua enumz.UA) (int64, bool) {
	rdb, ok := h.rdb(ctx)
	if !ok {
		return 0, false
	}
	prefix := incrUserTokenKeyPrefix()
	field := toUserTokenFactorField(svc, uid, ua)
	ttl := time.Duration(configz.UserRefreshTokenTTLs) * time.Second
	factor, e := cachez.IncrFactor(ctx, rdb, ttl, field, prefix, configz.UserTokenIntervalDays, 3)
	if !h.app.Check(ctx, e) {
		return 0, false
	}
	return factor, true
}

func (h *Cache) LoadUserTokenFactor(ctx context.Context, svc typez.Svc, uid uint64, ua enumz.UA) (int64, bool) {
	rdb, ok := h.rdb(ctx)
	if !ok {
		return 0, false
	}
	prefix := incrUserTokenKeyPrefix()
	field := toUserTokenFactorField(svc, uid, ua)
	factor, e := cachez.LoadFactor(ctx, rdb, field, prefix, configz.UserTokenIntervalDays)
	if !h.app.Check(ctx, e) {
		return 0, false
	}
	return factor, true
}
