package metrics

import "github.com/prometheus/client_golang/prometheus"

// UserOnlineCounter TODO - implement scheduler for check online users
var UserOnlineCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "user_online_total",
		Help: "Total number of user logins",
	},
)

func init() {
	Registry.MustRegister(UserOnlineCounter)
}
