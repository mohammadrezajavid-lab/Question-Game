package metricsserver

import (
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.project/go-fundamentals/gameapp/logger"
	"golang.project/go-fundamentals/gameapp/metrics"
	"net/http"
	"sync/atomic"
	"time"
)

type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode   int
	responseSize int64
}

func (w *responseWriterInterceptor) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterInterceptor) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.responseSize += int64(size)
	return size, err
}

type MetricsServer struct {
	config Config
	Server *http.Server
}

func NewMetricsServer(cfg Config) *MetricsServer {

	//metrics.Registry.MustRegister(metrics.HttpRequestCounter)
	handler := promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{})

	return &MetricsServer{
		config: cfg,
		Server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
			Handler: logRequest(handler),
		},
	}
}

func (ms *MetricsServer) Serve() {

	logger.Info(fmt.Sprintf("Starting metrics server on %s", ms.Server.Addr))

	if err := ms.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err, "Metrics server failed to start")
	}
}

func logRequest(next http.Handler) http.Handler {
	var requestCounter int64 = 0

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()

		interceptor := &responseWriterInterceptor{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(interceptor, r)

		latency := time.Since(startTime)
		currentRequestId := atomic.AddInt64(&requestCounter, 1)

		fields := []zapcore.Field{
			zap.Int64("request_id", currentRequestId),
			zap.String("protocol", r.Proto),
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", interceptor.statusCode),
			zap.String("latency", latency.String()),
			zap.String("host", r.Host),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int64("content_length", r.ContentLength),
			zap.Int64("response_size", interceptor.responseSize),
			zap.String("user_agent", r.UserAgent()),
		}

		zap.L().Named(logger.GetPackageFuncName(1)).Info("metrics", fields...)

	})
}
