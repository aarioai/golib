package metrics

import (
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/aarioai/airis/aa/httpsvr/response"
	"github.com/kataras/iris/v12"
	"net/http"
)

// 用于检测能否正常转发post的json数据
func postPing(ictx iris.Context) {
	defer ictx.Next()
	r, resp, _ := httpsvr.New(ictx)
	defer resp.CloseWith(r)
	ts, e0 := r.BodyInt64("timestamp")
	if e := resp.FirstError(e0); e != nil {
		return
	}
	resp.Write(map[string]int64{"timestamp": ts})
}

func (p *Party) WithPing() *Party {
	p.p.Head("/ping", response.StatusHandler(http.StatusOK))
	p.p.Get("/ping", response.WriteHandler("PONG"))
	p.p.Post("/ping", postPing)
	return p
}
