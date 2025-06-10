package metrics

import "github.com/prometheus/client_golang/prometheus"

var FailedZRemRedisCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_zrem_redis",
		Help: "Total number of failed ZRem redis",
	},
)

var FailedPublishedEventCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_published_event",
		Help: "Total number of failed published event",
	},
)

var PublishedEventCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "published_event",
		Help: "Total number of succeed published event",
	},
)

func init() {
	Registry.MustRegister(
		FailedZRemRedisCounter,
		FailedPublishedEventCounter,
		PublishedEventCounter,
	)
}
