package metrics

import "github.com/prometheus/client_golang/prometheus"

var FailedGetRandomQuestions = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_getrandomquestions_total",
		Help: "Total number of failed GetRandomQuestions",
	},
)

func init() {
	Registry.MustRegister(FailedGetRandomQuestions)
}
