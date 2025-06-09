package middleware

import (
	"errors"
	"github.com/labstack/echo/v4"
	"golang.project/go-fundamentals/gameapp/metrics"
	"net/http"
	"strconv"
)

func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// for run next middleware chain
			err := next(c)

			var statusCode int
			if err != nil {
				var httpError *echo.HTTPError
				if errors.As(err, &httpError) {
					statusCode = httpError.Code
				} else {
					statusCode = http.StatusInternalServerError
				}
			} else {
				statusCode = c.Response().Status
			}

			requestURI := c.Path()
			method := c.Request().Method

			// افزایش کانتر سفارشی شما
			metrics.HttpRequestCounter.WithLabelValues(strconv.Itoa(statusCode), requestURI, method).Inc()

			return err
		}
	}

	/*return echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{

		AfterNext: func(c echo.Context, err error) {
			statusCode := c.Response().Status

			if err != nil {
				if httpErr, ok := err.(*echo.HTTPError); ok && statusCode == http.StatusOK {
					statusCode = httpErr.Code
				}
			}

			requestURI := c.Path()
			method := c.Request().Method

			metrics.HttpRequestCounter.WithLabelValues(strconv.Itoa(statusCode), requestURI, method).Inc()
		},
		Registerer: metrics.Registry,
	})*/

}
