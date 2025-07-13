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
		Name: "failed_published_event_total",
		Help: "Total number of failed published event",
	},
)

var PublishedEventCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "published_event_total",
		Help: "Total number of succeed published event",
	}, []string{"event_name"},
)

var FailedSubscribeTopicCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "failed_subscribe_topic",
		Help: "Total number of failed subscribe topic",
	}, []string{"topic"},
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
		FailedSubscribeTopicCounter,
		RedisRequestsCounter,
	)
}
