package metrics

import "github.com/prometheus/client_golang/prometheus"

var FailedCreateSchCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_create_scheduler",
		Help: "Total number of failed create scheduler",
	},
)

var StartSchCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "start_scheduler_counter",
		Help: "Total number of start scheduler",
	},
)

var FailedCreateMatchWaitedUserJobCounter = prometheus.NewCounter(

	prometheus.CounterOpts{
		Name: "failed_create_match_waited_user_job",
		Help: "Total number of failed create match waited user job",
	},
)

var CreateMatchWaitedUserJobCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "create_match_waited_user_job",
		Help: "Total number of create match waited user job",
	},
)

var MatchWaitedUserFailedJobCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "match_waited_user_failed_job",
		Help: "Total number of MatchWaitedUser failed job",
	},
)

var MatchWaitedUserRunSuccessfullyJobCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "match_waited_user_run_successfully_job",
		Help: "Total number of run successfully MatchWaitedUser job",
	},
)

func init() {
	Registry.MustRegister(
		FailedCreateSchCounter,
		StartSchCounter,
		FailedCreateMatchWaitedUserJobCounter,
		CreateMatchWaitedUserJobCounter,
		MatchWaitedUserFailedJobCounter,
		MatchWaitedUserRunSuccessfullyJobCounter,
	)
}
