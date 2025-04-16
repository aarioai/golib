package cache

import (
	"context"
	"fmt"
	"github.com/aarioai/airis/aa/aenum"
	"github.com/aarioai/airis/pkg/types"
	"github.com/aarioai/golib/sdk/sms/config"
	"time"
)

// 短信、邮件验证码有效期不同，所以要独立出来
func (h *Cache) vericodeKey(pseudoId string, cn aenum.Country, phone string) string {
	cs := types.FormatUint(cn)
	account := cs + config.CacheDelimiter + phone + config.CacheDelimiter + pseudoId
	return fmt.Sprintf(config.CacheVericodeKeyFormat, account)
}

func (h *Cache) vericodeLimitKey(cn aenum.Country, phone string) string {
	cs := types.FormatUint(cn)
	account := cs + config.CacheDelimiter + phone
	return fmt.Sprintf(config.CacheVericodeLimitKeyFormat, account)
}

// ApplySmsVericodeSendingPermission 短信、邮件验证码有效期不同，所以要独立出来
// 统一验证码发送，验证码对稳定性要求比较高
// 各个子业务，可以自行决定是否使用自己的短信验证码
// 未登录短信验证码，需要根据account来，故不能用通用的限流
func (h *Cache) ApplySmsVericodeSendingPermission(ctx context.Context, cn aenum.Country, phone string, periodLimit int) (time.Duration, bool) {
	rdb, ok := h.rdb(ctx)
	if !ok {
		return 0, false
	}

	k := h.vericodeLimitKey(cn, phone)
	maxTTL := config.VericodePeriodTTL

	ttl, err := rdb.TTL(ctx, k).Result()
	if err != nil || ttl.Seconds() < 1.0 {
		_, err = rdb.SetEx(ctx, k, 1, maxTTL).Result()
		return maxTTL, h.app.CheckErrors(ctx, err)
	}
	limit, _ := rdb.Incr(ctx, k).Result()
	if limit < 1 {
		_, err = rdb.SetEx(ctx, k, 1, maxTTL).Result()
		return maxTTL, h.app.CheckErrors(ctx, err)
	}
	return ttl, int64(periodLimit) <= limit
}

// LoadSmsVericode 短信发送存在延时性问题，每10分钟内，重复发送相同的验证码
func (h *Cache) LoadSmsVericode(ctx context.Context, cn aenum.Country, phone, pseudoId string) (string, bool) {
	k := h.vericodeKey(pseudoId, cn, phone)
	rdb, ok := h.rdb(ctx)
	if !ok {
		return "", false
	}

	ttl, err := rdb.TTL(ctx, k).Result()
	// ttl 必须要大于30秒，才能重复使用
	if err != nil || ttl.Seconds() < 30.0 {
		return "", false
	}
	var vericode string
	vericode, err = rdb.Get(ctx, k).Result()
	if !h.app.CheckErrors(ctx, err) {
		return "", false
	}
	return vericode, true
}

func (h *Cache) LoadAndDeleteVericode(ctx context.Context, cn aenum.Country, phone, pseudoId string) (string, bool) {
	k := h.vericodeKey(pseudoId, cn, phone)
	rdb, ok := h.rdb(ctx)
	if !ok {
		return "", false
	}

	// 这里一定要用 defer ，不然先删除，后面那行就找不到了
	defer rdb.Del(ctx, k)
	vericode, err := rdb.Get(ctx, k).Result()
	return vericode, h.app.CheckErrors(ctx, err)
}

func (h *Cache) CacheSmsVericode(ctx context.Context, cn aenum.Country, phone, vericode, pseudoId string) bool {
	rdb, ok := h.rdb(ctx)
	if !ok {
		return false
	}

	k := h.vericodeKey(pseudoId, cn, phone)
	ttl := config.VericodeTTL
	_, err := rdb.SetEx(ctx, k, vericode, ttl).Result()
	return h.app.CheckErrors(ctx, err)
}
