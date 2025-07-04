package metrics

import "github.com/prometheus/client_golang/prometheus"

var UserOnlineCounter = prometheus.NewGauge(
	prometheus.GaugeOpts{
		Name: "user_online_total",
		Help: "Total number of user logins",
	},
)

func init() {
	Registry.MustRegister(UserOnlineCounter)
}
