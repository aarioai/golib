package metrics

import (
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

func WithMetrics(p iris.Party) iris.Party {
	p.Use(prometheusMiddleware)
	p.Get("/metrics", iris.FromStd(promhttp.Handler()))
	return p
}
