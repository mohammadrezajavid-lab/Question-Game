package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"runtime"
	"time"
)

var GoActiveGoroutinesGauge = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "go_active_goroutines",
		Help: "Total number of active goroutines",
	},
)

var GoActiveGoroutinesServiceGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "service_goroutines",
		Help: "Total number of active goroutines per service",
	}, []string{"service"},
)

func init() {
	Registry.MustRegister(
		GoActiveGoroutinesGauge,
		GoActiveGoroutinesServiceGauge,
	)

	go func() {
		for {
			GoActiveGoroutinesGauge.Set(float64(runtime.NumGoroutine()))
			time.Sleep(5 * time.Second)
		}
	}()
}
