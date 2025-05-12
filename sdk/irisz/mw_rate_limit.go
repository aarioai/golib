package irisz

import (
	"fmt"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/aconfig"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/types"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/rate"
	"strconv"
	"strings"
)

type RateLimitType string

const (
	RateLimitAPI  = "api"
	RateLimitView = "view"
)

func parseDurationConfig(c *aconfig.Config, keys []string) string {
	for _, key := range keys {
		v := c.GetString(key)
		if v != "" {
			return strings.ReplaceAll(v, " ", "")
		}
	}
	return ""
}
func RateLimitMiddleware(p iris.Party, app *aa.App, sectionName string, t RateLimitType) iris.Party {
	keys := []string{
		fmt.Sprintf("%s.%s_rate_limit", sectionName, t),
		fmt.Sprintf("%s.rate_limit", sectionName),
		fmt.Sprintf("app.%s_rate_limit", sectionName),
		"app.rate_limit",
	}
	v := parseDurationConfig(app.Config, keys)
	if v == "" {
		return p
	}

	limits := strings.Split(v, ",")
	limit, err := types.Atoi(limits[0])
	ae.PanicOnErrors(err)
	burst := limit // 默认桶的容量等于每秒消耗最高量，如果一直没有消费掉，则持续往桶里增加，超过桶上限的令牌丢弃
	if len(limits) > 0 && limits[1] != "" {
		burst, err = strconv.Atoi(limits[1])
		ae.PanicOnErrors(err)
	}
	// 限流 Limit(limit 每秒放token数，即QPS, burst 令牌桶大小，即最大并发数)
	// limit = 1000 = rate.Every(1 * time.Millisecond)   每毫秒放1个，即每秒1000个
	// Limit(500,1000) ==> 每秒500个，桶的容量为1000。如果一直没有消费掉，则持续往桶里增加，超过桶上限的令牌丢弃
	//   如果请求时桶里没有令牌，则被限流
	// E.g. Limit(1, 5)
	// E.g. Limit(..., PurgeEvery(time.Minute, 5*time.Minute)) to check every 1 minute if a client's last visit was 5 minutes ago ("old" entry) and remove it from the memory.
	if len(limits) == 4 {
		every := types.ParseDuration(limits[2])
		maxLifetime := types.ParseDuration(limits[3])
		p.Use(rate.Limit(float64(limit), burst, rate.PurgeEvery(every, maxLifetime)))
		return p
	}
	p.Use(rate.Limit(float64(limit), burst))
	return p
}
