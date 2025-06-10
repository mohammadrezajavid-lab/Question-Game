package metrics

import "github.com/prometheus/client_golang/prometheus"

var FailedUpsertPresenceCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_upsert_presence_total",
		Help: "Total number of failed upsert presence service",
	},
)

var FailedGetPresenceServiceCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_get_presence_service_total",
		Help: "Total number of failed GetPresence service",
	},
)

func init() {
	Registry.MustRegister(
		FailedUpsertPresenceCounter,
		FailedGetPresenceServiceCounter,
	)
}
