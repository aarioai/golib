package cache

import (
	"context"
	"github.com/aarioai/golib/sdk/auth/configz"
	"github.com/redis/go-redis/v9"
)

// 指纹是否生成，KV结构

func fingerprintKey(uuid string) string { return configz.CachePrefix + "mmc:fp:uuid:" + uuid }
func fingerprintBFKey() string          { return configz.CachePrefix + "mmc:fp:bf:" }

/*
	bloom filter 算法只适合可以反复重复使用的数据查找，redis 带这个module；bloom filter 其实就是多次hash。算法更适合做黑名单、白名单。

这里uuid 是临时性的，不适合用 bloom filter
https://redis.io/commands/?name=bf.
Bloom vs. Cuckoo filters
Bloom filters typically exhibit better performance and scalability when inserting items (so if you're often adding items to your dataset, then a Bloom filter may be ideal). Cuckoo filters are quicker on check operations and also allow deletions.

	Bloom filters 不能删除之前插入过的数据（多hash原理：https://www.icode9.com/content-4-1074321.html）
*/
func (h *Cache) CacheFingerprintUUID(ctx context.Context, uuid string) bool {
	rdb, ok := h.rdb(ctx)
	if !ok {
		return false
	}

	k := fingerprintKey(uuid)
	err := rdb.SetEx(ctx, k, 1, configz.MmcFingerprintTTL).Err()
	return h.app.CheckErrors(ctx, err)
}

// fingerprint 使用可以持续一段时间
func (h *Cache) CheckMmcFingerprintUUID(ctx context.Context, uuid string) bool {
	rdb, ok := h.rdb(ctx)
	if !ok {
		return true // redis错误，就放开
	}

	k := fingerprintKey(uuid)
	ttl, err := rdb.TTL(ctx, k).Result()
	if err != nil || ttl.Seconds() < 1.0 {
		return false
	}
	return true
}

func (h *Cache) bfAddFingerprint(ctx context.Context, hash string) (exists, ok bool) {
	var rdb *redis.Client
	rdb, ok = h.rdb(ctx)
	if !ok {
		return false, false
	}
	inserted, err := rdb.Do(ctx, "BF.ADD", fingerprintBFKey(), hash).Bool()
	if h.app.CheckErrors(ctx, err) {
		return !inserted, true
	}
	return false, false
}
