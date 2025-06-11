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

var RedisRequestsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "redis_requests_total",
		Help: "How many Redis requests processed, partitioned by status",
	}, []string{"status"},
)

func init() {
	Registry.MustRegister(
		FailedZRemRedisCounter,
		FailedPublishedEventCounter,
		PublishedEventCounter,
		RedisRequestsCounter,
	)
}
