package irisz

import (
	"github.com/aarioai/airis/aa/httpsvr"
	"github.com/aarioai/airis/aa/httpsvr/response"
	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

var (
	totalRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "Total number of HTTP requests",
	}, []string{"path", "method", "status"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_duration_seconds",
		Help:    "Duration of HTTP requests",
		Buckets: prometheus.DefBuckets,
	}, []string{"path", "method", "status"})
)

func prometheusMiddleware(ictx iris.Context) {
	start := time.Now()
	path := ictx.Path()
	ictx.Next()
	duration := time.Since(start).Seconds()
	status := ictx.GetStatusCode()

	totalRequests.WithLabelValues(ictx.Method(), path, http.StatusText(status)).Inc()
	requestDuration.WithLabelValues(ictx.Method(), path, http.StatusText(status)).Observe(duration)
}

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

func (w *Irisz) RouteMetrics(p iris.Party) {
	p.Use(prometheusMiddleware)
	p.Head("/ping", response.StatusHandler(http.StatusOK))
	p.Get("/ping", response.WriteHandler("PONG"))
	p.Post("/ping", postPing)
	p.Get("/metrics", iris.FromStd(promhttp.Handler()))
}
