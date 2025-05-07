package adapter

import (
	"errors"
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis-driver/driver/mysqli"
	"github.com/aarioai/airis/aa"
	"github.com/aarioai/airis/aa/ae"
	"github.com/aarioai/airis/pkg/afmt"
	"github.com/casbin/casbin/v2/persist"
)

type Adapter struct {
	app          *aa.App
	mysqlSection string
}

func New(app *aa.App, mysqlSection string) persist.Adapter {
	return &Adapter{
		app:          app,
		mysqlSection: mysqlSection,
	}
}

func (a *Adapter) db() *mysqli.DB {
	return mysqli.NewDriver(driver.NewMysqlPool(a.app, a.mysqlSection))
}

func NewAdapterError(msg string, args ...any) error {
	return errors.New(afmt.Sprintf("sdk_casbinz: adapter "+msg, args...))
}

func handleDriverError(e *ae.Error) error {
	if e == nil || e.IsNotFound() {
		return nil
	}
	return errors.New("sdk_casbinz: adapter " + e.Text())
}
