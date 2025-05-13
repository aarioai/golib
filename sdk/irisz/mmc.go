package irisz

import (
	"github.com/aarioai/golib/enumz"
	"github.com/kataras/iris/v12"
)

func SetFingerprintServerTime(ictx iris.Context, ms int64) {
	ictx.Values().Set(enumz.IctxParamFingerprintServerTime, ms)
}
func FingerprintServerTime(ictx iris.Context) int64 {
	ms, _ := ictx.Values().GetInt64(enumz.IctxParamFingerprintServerTime)
	return ms
}
