package adapter

import (
	"errors"
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis-driver/driver/mysqli"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
)

const prefix = "libsdk_casbinz:adapter: "

type Adapter struct {
	app          *aa.App
	mysqlSection string
}

// New extends persist.Adapter
func New(app *aa.App, mysqlSection string) *Adapter {
	return &Adapter{
		app:          app,
		mysqlSection: mysqlSection,
	}
}

func (a *Adapter) db() *mysqli.DB {
	return mysqli.NewDriver(driver.NewMysqlPool(a.app, a.mysqlSection))
}

func NewAdapterError(msg string, args ...any) error {
	return errors.New(afmt.Sprintf(prefix+"adapter "+msg, args...))
}

func handleDriverError(e *ae.Error) error {
	if e == nil || e.IsNotFound() {
		return nil
	}
	return errors.New(prefix + "adapter " + e.Text())
}
