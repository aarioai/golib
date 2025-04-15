package alg

import (
	"math"
	"time"
)

// ranking 必须要加上律师更新时的时间戳差，这样才能活跃的律师能更好展示
type Ranking int64                       // uint24
const RankingTimestampStart = 1694999880 // 2023-09-18 09:18:00

// 对律所、律师、商品等排序算法
const ActivePointMax = 65535

func ToRanking(v int64, now time.Time) Ranking {
	return Ranking(now.Unix() - RankingTimestampStart)
}

// rang 有可能是 [-100,0] 比如点击不喜欢

func TryRank(ok bool, min, max int64, p float64) int64 {
	if !ok || max <= min {
		return 0
	}
	a := max - min
	if p != 1.0 {
		a = int64(math.Ceil(p * float64(a)))
	}
	return a + min
}

/*
		基于 Wilson Interval 算法  范围 -2147483648~2147483647
	 转为：uint64  ->  42949 67295  以下的，都是复数

		18446744073709551615
*/
