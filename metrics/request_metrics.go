package metrics

import "github.com/prometheus/client_golang/prometheus"

var HttpRequestCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests received",
	}, []string{"status", "path", "method"},
)

func init() {
	Registry.MustRegister(HttpRequestCounter)
}
