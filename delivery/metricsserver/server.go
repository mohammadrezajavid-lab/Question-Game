package metricsserver

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metricsMiddleware "golang.project/go-fundamentals/gameapp/delivery/metricsserver/middleware"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"net/http"
)

type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type MetricsServer struct {
	config Config
	Server *http.Server
}

func NewMetricsServer(cfg Config) *MetricsServer {

	handler := promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{})

	return &MetricsServer{
		config: cfg,
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler: metricsMiddleware.LogRequestMiddleware(handler),
		},
	}
}

func (ms *MetricsServer) Serve() {

	logger.Info(fmt.Sprintf("Starting metrics server on %s", ms.Server.Addr))

	if err := ms.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err, "Metrics server failed to start")
	}
}
