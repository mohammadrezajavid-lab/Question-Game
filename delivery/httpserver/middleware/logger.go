package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ZapLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError:      true,
		LogStatus:        true,
		LogLatency:       true,
		LogProtocol:      true,
		LogHost:          true,
		LogRemoteIP:      true,
		LogMethod:        true,
		LogURI:           true,
		LogRequestID:     true,
		LogUserAgent:     true,
		LogContentLength: true,
		LogResponseSize:  true,
		LogError:         true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {

			fields := []zapcore.Field{
				zap.String("request_id", v.RequestID),
				zap.String("protocol", v.Protocol),
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
				zap.String("latency", v.Latency.String()),
				zap.String("host", v.Host),
				zap.String("remote_ip", v.RemoteIP),
				zap.String("content_length", v.ContentLength),
				zap.Int64("response_size", v.ResponseSize),
				zap.String("user_agent", v.UserAgent),
			}

			statusCode := v.Status
			if v.Error == nil {
				switch {
				case statusCode >= 500:
					zap.L().Error("internal server error", fields...)
				case statusCode >= 400:
					zap.L().Warn("client warning", fields...)
				case statusCode >= 300:
					zap.L().Info("redirection", fields...)
				default:
					zap.L().Info("success", fields...)
				}
			} else {
				zap.L().Error("request error", fields...)
			}
			return nil
		},
	})
}
