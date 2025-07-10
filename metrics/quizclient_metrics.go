package metrics

import "github.com/prometheus/client_golang/prometheus"

var FailedOpenQuizClientGRPCConnCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_open_quiz_client_grpc_connection",
		Help: "Total number of failed open QuizClient GRPC connection",
	},
)

func init() {
	Registry.MustRegister(FailedOpenQuizClientGRPCConnCounter)
}
