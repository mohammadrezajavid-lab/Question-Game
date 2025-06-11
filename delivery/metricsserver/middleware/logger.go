package middleware

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.project/go-fundamentals/gameapp/logger"
	"net/http"
	"sync/atomic"
	"time"
)

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

func LogRequestMiddleware(next http.Handler) http.Handler {
	var requestCounter int64 = 0

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		startTime := time.Now()

		responseWriter := &responseWriterInterceptor{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(responseWriter, r)

		latency := time.Since(startTime)
		currentRequestId := atomic.AddInt64(&requestCounter, 1)

		fields := []zapcore.Field{
			zap.Int64("request_id", currentRequestId),
			zap.String("protocol", r.Proto),
			zap.String("method", r.Method),
			zap.String("uri", r.RequestURI),
			zap.Int("status", responseWriter.statusCode),
			zap.String("latency", latency.String()),
			zap.String("host", r.Host),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int64("content_length", r.ContentLength),
			zap.Int64("response_size", responseWriter.responseSize),
			zap.String("user_agent", r.UserAgent()),
		}

		zap.L().Named(logger.GetPackageFuncName(1)).Info("metrics", fields...)

	})
}
