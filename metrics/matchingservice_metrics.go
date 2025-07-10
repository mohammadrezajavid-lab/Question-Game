package metrics

import "github.com/prometheus/client_golang/prometheus"

var FailedGetPresenceClientCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_get_presence_client_total",
		Help: "Total number of failed GetPresence grpc client",
	},
)

var FailedAddToWaitingListCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_add_to_waiting_list_total",
		Help: "Total number of failed AddToWaitingList",
	},
)

var FailedGetWaitedUsersByCategoryCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_getWaitedUsersByCategory_total",
		Help: "Total number of failed get user from WaitingList by category",
	},
)

func init() {
	Registry.MustRegister(
		FailedGetPresenceClientCounter,
		FailedAddToWaitingListCounter,
		FailedGetWaitedUsersByCategoryCounter,
	)
}
