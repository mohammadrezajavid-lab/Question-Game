package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"golang.project/go-fundamentals/gameapp/metrics"
	"net/http"
	"strconv"
)

func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			metrics.ActiveRequestsGauge.Inc()

			requestURI := c.Path()
			method := c.Request().Method

			timer := prometheus.NewTimer(metrics.HTTPLatency.WithLabelValues(requestURI))

			err := next(c)

			var statusCode int
			if err != nil {
				var httpError *echo.HTTPError
				// type assertion: err.(*echo.HTTPError)
				if errors.As(err, &httpError) {
					statusCode = httpError.Code
				} else {
					statusCode = http.StatusInternalServerError
				}
			} else {
				statusCode = c.Response().Status
			}
			
			metrics.ActiveRequestsGauge.Dec()
			metrics.HttpRequestCounter.WithLabelValues(strconv.Itoa(statusCode), requestURI, method).Inc()
			timer.ObserveDuration()

			return err
		}
	}
}
