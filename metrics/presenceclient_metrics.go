package metrics

import "github.com/prometheus/client_golang/prometheus"

var FailedOpenPresenceClientGRPCConnCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_open_presence_client_grpc_connection",
		Help: "Total number of failed open PresenceClient GRPC connection",
	},
)

func init() {
	Registry.MustRegister(FailedOpenPresenceClientGRPCConnCounter)
}
