package cachez

import (
	"context"
	"crypto/rand"
	"github.com/aarioai/airis-driver/driver/redishelper"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/types"
	"github.com/redis/go-redis/v9"
	"math/big"
	"time"
)

// Interval Increment Factor 具有时效的，递增数字

func incrFactorKeyName(prefix string, intervalDays uint8, prev bool) string {
	return redishelper.DailyKey(intervalDays, prev, func(n uint8) string {
		return prefix + types.FormatUint8(n)
	})
}

// IncrFactor 增加Incr
// field 字段，如 $uid:$ua
// ttl 时效
// prefix key 前缀
// intervalDays 日间隔
// incrMax 每次增加最大值，若<=1，则增加1；若>1，则随机一个 1~incrMax 之间整数值
func IncrFactor(ctx context.Context, rdb *redis.Client, ttl time.Duration, field, prefix string, intervalDays uint8, incrMax int64) (int64, *ae.Error) {
	k := incrFactorKeyName(prefix, intervalDays, false)

	// 随机增加一个数字
	if incrMax > 1 {
		maxIncr := big.NewInt(incrMax)
		incr, _ := rand.Int(rand.Reader, maxIncr)
		incrMax = incr.Int64()
	}
	if incrMax < 1 {
		incrMax = 1
	}

	factor, e := redishelper.HIncrBy(ctx, rdb, ttl, k, field, incrMax)
	if e != nil {
		return 0, e
	}
	if factor == 0 {
		return 0, NewE("hincrby %d to %s -> %s failed", incrMax, k, field)
	}
	return factor, nil
}

func LoadFactor(ctx context.Context, rdb *redis.Client, field, prefix string, intervalDays uint8) (int64, *ae.Error) {
	k := incrFactorKeyName(prefix, intervalDays, false)
	id, err := rdb.HGet(ctx, k, field).Int64()
	if err == nil {
		return id, nil
	}
	// 查询上一周期
	k = incrFactorKeyName(prefix, intervalDays, true)
	id, err = rdb.HGet(ctx, k, field).Int64()
	if err == nil {
		return id, nil
	}
	return 0, ae.ErrorNotFound
}
