package metrics

import "github.com/prometheus/client_golang/prometheus"

var FailedReadPompCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_readpomp_total",
		Help: "Total number of failed ReadPomp in websocket",
	},
)

var FailedWritePompCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_writepomp_total",
		Help: "Total number of failed WritePomp in websocket",
	},
)

var InvalidHeartBeatMessageCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "invalid_heartbeat_message_total",
		Help: "Total number of invalid heartbeat message for any client",
	}, []string{"user_id"},
)

var ClosingConnectionWebsocket = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "close_connection_websocket_total",
		Help: "Total number of Close connection websocket",
	},
)

func init() {
	Registry.MustRegister(
		FailedReadPompCounter,
		FailedWritePompCounter,
		InvalidHeartBeatMessageCounter,
		ClosingConnectionWebsocket,
	)
}
