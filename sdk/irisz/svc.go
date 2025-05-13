package irisz

import (
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/golib/enumz"
	"github.com/kataras/iris/v12"
	"strconv"
)

func SetAg(ictx iris.Context, f string) bool {
	_, ok := ictx.Values().Set(enumz.ParamAg, f)
	return ok
}
func Ag(ictx iris.Context) (string, *ae.Error) {
	// 先尝试从 param 读取
	ag := ictx.URLParam(enumz.ParamAg)
	if ag == "" {
		ag = ictx.Values().GetString(enumz.ParamAg)
	}
	if ag == "" {
		return "", ae.NewBadParam(enumz.ParamAg)
	}
	return ag, nil
}

func SetFromUid(ictx iris.Context) (uint64, bool) {
	ag, _ := Ag(ictx)
	fromUid, _ := strconv.ParseUint(ag, 36, 64)
	return fromUid, fromUid > 0
}
