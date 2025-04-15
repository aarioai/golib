package coding

import (
	"sync/atomic"
	"time"
)

const (
	baseTime = 1330768499 // 基准时间   让最后总值从 1000000000-0000 开始，到 9999999999-9999  14位数字，可以表达到 2100年
	maxSeq   = 65536      // 65536 = 2 << 15
)

func incrDefaultUuidSeq(addr *uint64) uint64 {
	//atomic.CompareAndSwapUint64(addr, math.MaxUint64, 0)   // 数值足够大，考虑增加到最大值的可能性极低
	return atomic.AddUint64(addr, 1)
}

//	14位数字，支持1秒钟并发 65536 = 2 << 15
//
// @example    var atomicAddr uint64      codingUint64Id(&atomicAddr)
func Uint64Id(addr *uint64) uint64 {
	return NewUint64Id(time.Now(), addr, 15) // 65536 = 2 << 15
}
func NewUint64Id(t time.Time, addr *uint64, n int) uint64 {
	now := t.Unix() - baseTime
	return NewU64Id(now, addr, n)
}
func NewU64Id(tm int64, addr *uint64, n int) uint64 {
	seq := incrDefaultUuidSeq(addr)
	return ToUint64Id(tm, seq, n)
}

// 反转short uuid
func ParseUint64Id(id uint64, n int) (ts int64, seq uint64) {
	// 65536 = 2 << n
	seq = ((2 << n) - 1) & id
	ts = int64(id >> n)
	if ts < 0 {
		ts = 0
	}
	return
}

func ToUint64Id(ts int64, seq uint64, n int) uint64 {
	id := uint64(ts) << n
	id |= seq % (2 << n)
	return id
}
