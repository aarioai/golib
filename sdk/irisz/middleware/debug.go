package middleware

import (
	"github.com/aarioai/airis/aa/acontext"
	"github.com/aarioai/airis/aa/httpsvr/request"
	"github.com/aarioai/golib/enumz"
	"github.com/kataras/iris/v12"
	"os/exec"
	"time"
)

var devSyncFileAt time.Time

// SyncLocalFiles 每次请求刷新下文件，主要用于虚拟主机、宿主之间文件修改不同步问题
// @warn 应仅用于local使用虚拟机的情况下
func (w *Middleware) SyncLocalFiles(ictx iris.Context) {
	defer ictx.Next()
	if devSyncFileAt.Add(time.Second).After(time.Now()) {
		return // 至少隔1秒 才sync，避免太频繁
	}
	devSyncFileAt = time.Now()
	w.app.CheckErrors(acontext.FromIris(ictx), exec.Command("bash", "-c", "sync").Run())
}

func (w *Middleware) WithDebugToken(token string) *Middleware {
	w.debugToken = token
	return w
}

func (w *Middleware) WithMockToken(token string) *Middleware {
	w.mockToken = token
	return w
}

// LoadDebugTokens mock 需要同时后台开启，并且客户端cookie或GET 中 _mock=1
func (w *Middleware) LoadDebugTokens(ictx iris.Context) {
	defer ictx.Next()

	if w.debugToken != "" && request.HeaderWild(ictx, enumz.HeaderDebugToken) == w.debugToken {
		ictx.Values().Set(enumz.IctxClientDebug, true)
	}
	if w.mockToken != "" && request.HeaderWild(ictx, enumz.HeaderMockToken) == w.mockToken {
		ictx.Values().Set(enumz.IctxClientMock, true)
	}
}
