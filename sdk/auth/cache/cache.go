package cache

import (
	"context"
	"github.com/aarioai/airis-driver/driver"
	"github.com/aarioai/airis/aa"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache struct {
	app     *aa.App
	loc     *time.Location
	section string
}

func New(app *aa.App, section string) *Cache {
	return &Cache{
		app:     app,
		loc:     app.Config.TimeLocation,
		section: section,
	}
}

// go redis will create connection pool automatically. so its no need to close a connection.
// @doc https://pkg.go.dev/github.com/redis/go-redis/v9#Client.Close
// It is rare to Close a Client, as the Client is meant to be long-lived and shared between many goroutines.
func (h *Cache) rdb(ctx context.Context) (*redis.Client, bool) {
	cli, e := driver.NewRedisPool(h.app, h.section)
	if e != nil {
		h.app.Log.Error(ctx, e.Text())
		return nil, false
	}
	return cli, true
}
