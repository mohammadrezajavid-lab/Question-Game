package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var HttpRequestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests received",
	}, []string{"status", "path", "method"},
)

var ActiveRequestsGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "http_active_requests",
		Help: "Number of active connections to the service",
	},
)

var HTTPLatency = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "http_latency_time_secondes",
		Help: "Latency of HTTP requests",
	}, []string{"path"},
)

func init() {
	Registry.MustRegister(
		HttpRequestCounter,
		ActiveRequestsGauge,
		HTTPLatency,
	)
}
