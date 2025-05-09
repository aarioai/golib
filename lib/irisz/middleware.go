package irisz

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"time"
)

type Middleware struct {
	app              *aa.App
	loc              *time.Location
	authRedisSection string
}

func New(app *aa.App, authRedisSection string) *Middleware {
	return &Middleware{app: app,
		loc:              app.Config.TimeLocation,
		authRedisSection: authRedisSection,
	}
}

func NewCode(code int, format string, args ...any) *ae.Error {
	return ae.New(code, afmt.Sprintf("lib_irisz: "+format, args...))
}

func NewE(format string, args ...any) *ae.Error {
	return ae.NewE("lib_irisz: "+format, args...)
}

func NewError(err error) *ae.Error {
	if err == nil {
		return nil
	}
	return ae.NewE("lib_irisz: " + err.Error())
}
