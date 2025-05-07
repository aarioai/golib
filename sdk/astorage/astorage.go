package astorage

import (
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis-driver/driver/mysqli"
	"github.com/aarioai/airis/aa"
)

type AStorage struct {
	app             *aa.App
	mysqlCfgSection string
}

func New(app *aa.App, mysqlCfgSection string) *AStorage {
	return &AStorage{
		app:             app,
		mysqlCfgSection: mysqlCfgSection,
	}
}

func (g *AStorage) db() *mysqli.DB {
	return mysqli.NewDriver(driver.NewMysql(g.app, g.mysqlCfgSection))
}
