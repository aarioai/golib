package middleware

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/aarioai/golib/typez"
	"github.com/kataras/iris/v12"
	"time"
)

const prefix = "libsdk_auth: "

type Middleware struct {
	app           *aa.App
	loc           *time.Location
	configSection string
	parser        func(iris.Context) (svc typez.Svc, uid, vuid uint64, ttl int64, atoken string, e *ae.Error)
	tidyParser    bool
}

// New 仅限于共享Redis内存的服务使用，否则应另写接口使用GRPC方式调用
func New(app *aa.App) *Middleware {
	return &Middleware{
		app: app,
		loc: app.Config.TimeLocation,
	}
}

func (w *Middleware) WithRedis(configSection string) *Middleware {
	w.configSection = configSection
	return w
}

func (w *Middleware) WithParser(parser func(iris.Context) (svc typez.Svc, uid, vuid uint64, ttl int64, atoken string, e *ae.Error), tidy bool) *Middleware {
	w.parser = parser
	w.tidyParser = tidy
	return w
}

func NewCode(code int, format string, args ...any) *ae.Error {
	return ae.New(code, afmt.Sprintf(prefix+format, args...))
}

func NewE(format string, args ...any) *ae.Error {
	return ae.NewE(prefix+format, args...)
}

func NewError(err error) *ae.Error {
	if err == nil {
		return nil
	}
	return ae.NewE(prefix + err.Error())
}
