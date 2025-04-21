package cache

import (
	"context"
	"github.com/aarioai/golib/sdk/auth/configz"
	"time"
)

func apolloUpdateAtKey() string {
	return configz.CachePrefix + "apollo:upat"
}

func (h *Cache) CacheApolloUpdatedAt(ctx context.Context, updatedAt int64) bool {
	rdb, ok := h.rdb(ctx)
	if !ok {
		return false
	}
	_, err := rdb.Set(ctx, apolloUpdateAtKey(), updatedAt, time.Hour*24*365).Result()
	return h.app.CheckErrors(ctx, err)
}

func (h *Cache) LoadApolloUpdatedAt(ctx context.Context) (int64, bool) {
	rdb, ok := h.rdb(ctx)
	if !ok {
		return 0, false
	}
	authAt, err := rdb.Get(ctx, apolloUpdateAtKey()).Int64()
	return authAt, h.app.CheckErrors(ctx, err)
}
