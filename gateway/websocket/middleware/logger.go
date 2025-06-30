package middleware

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.project/go-fundamentals/gameapp/logger"
	"net/http"
	"sync/atomic"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	var requestCounter int64 = 0

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			currentRequestId := atomic.AddInt64(&requestCounter, 1)

			fields := []zapcore.Field{
				zap.Int64("request_id", currentRequestId),
				zap.String("protocol", r.Proto),
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.String("host", r.Host),
				zap.String("remote_addr", r.RemoteAddr),
				zap.Int64("content_length", r.ContentLength),
				zap.String("user_agent", r.UserAgent()),
			}

			zap.L().Named(logger.GetPackageFuncName(1)).Info("websocket", fields...)

			next.ServeHTTP(w, r)
		},
	)
}
