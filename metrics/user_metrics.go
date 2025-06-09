package metrics

import "github.com/prometheus/client_golang/prometheus"

var UserLoginCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "user_login_total",
		Help: "Total number of user logins",
	},
)

func init() {
	Registry.MustRegister(UserLoginCounter)
}
