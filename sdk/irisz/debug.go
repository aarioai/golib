package irisz

import (
	"github.com/aarioai/golib/enumz"
	"github.com/kataras/iris/v12"
)

func ClientDebug(ictx iris.Context) bool {
	return ictx.Values().GetBoolDefault(enumz.IctxClientDebug, false)
}

func ClientMock(ictx iris.Context) bool {
	return ictx.Values().GetBoolDefault(enumz.IctxClientMock, false)
}
