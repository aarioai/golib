package irisz

import (
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"time"
)

const prefix = "libsdk_irisz: "

type Irisz struct {
	app              *aa.App
	loc              *time.Location
	authRedisSection string
}

func New(app *aa.App) *Irisz {
	return &Irisz{
		app: app,
		loc: app.Config.TimeLocation,
	}
}
func NewError(err error) *ae.Error {
	if err == nil {
		return nil
	}
	return ae.NewE(prefix + err.Error())
}

func PanicOnErrors(errs ...error) {
	for _, err := range errs {
		if err != nil {
			panic(prefix + err.Error())
		}
	}
}
