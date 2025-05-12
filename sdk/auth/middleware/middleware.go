package middleware

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"time"
)

const prefix = "libsdk_auth: "

type Middleware struct {
	app              *aa.App
	loc              *time.Location
	authRedisSection string
}

// New 仅限于共享Redis内存的服务使用，否则应另写接口使用GRPC方式调用
func New(app *aa.App, authRedisSection string) *Middleware {
	return &Middleware{app: app,
		loc:              app.Config.TimeLocation,
		authRedisSection: authRedisSection,
	}
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
