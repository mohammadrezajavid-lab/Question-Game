package metrics

import "github.com/prometheus/client_golang/prometheus"

var FailedGetQuestionsCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_getquestion_total",
		Help: "Total number of failed GetQuestion from Database",
	},
)

var FailedStoreQuiz = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_storequiz_total",
		Help: "Total number of failed store quiz in key-value storage",
	},
)

var FailedCreateGame = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_creategame_total",
		Help: "Total number of failed CreateGame",
	},
)

var SuccesGameCreated = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "succes_gamecreate_total",
		Help: "Total number of success CreateGame",
	},
)

var FailedFetchQuizFromPool = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_getquiz_total",
		Help: "Total number of failed fetch quiz from pool",
	},
)

var FailedCreatePlayers = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_createplayers_total",
		Help: "Total number of failed CreatePlayers",
	},
)

func init() {
	Registry.MustRegister(
		FailedGetQuestionsCounter,
		FailedStoreQuiz,
		FailedCreateGame,
		SuccesGameCreated,
		FailedFetchQuizFromPool,
		FailedCreatePlayers,
	)
}
