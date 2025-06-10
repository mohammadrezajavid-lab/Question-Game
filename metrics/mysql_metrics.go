package metrics

import "github.com/prometheus/client_golang/prometheus"

var FailedOpenMySQLConnCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "failed_open_mysql_connection",
		Help: "Total number of failed open mysql connection",
	},
)

var ConnectionRefusedMySQLCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "connection_refused_mysql",
		Help: "Total number of Connection refused MySQL",
	},
)

var DBQueryCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "db_queries_total",
		Help: "Total number of queries to MySQL database",
	}, []string{"query_type"},
)

var DBFailedQueryCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "db_failed_queries_total",
		Help: "Total number of queries to MySQL database is failed",
	}, []string{"query_type"},
)

func init() {
	Registry.MustRegister(
		FailedOpenMySQLConnCounter,
		ConnectionRefusedMySQLCounter,
		DBQueryCounter,
		DBFailedQueryCounter,
	)
}
