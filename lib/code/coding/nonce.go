package coding

import (
	"strconv"
	"time"
)

// TimeNonce
// @TODO 未来优化
func TimeNonce() (string, int64) {
	ts := time.Now().Unix()
	nonce := strconv.FormatInt(ts, 10)
	return nonce, ts
}
