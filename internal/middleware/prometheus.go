package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/rajabhishekmaurya/ecom/libs/monitoring"
)

func PrometheusMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			duration := time.Since(start).Seconds()

			monitoring.HTTPRequestTotal.
				WithLabelValues(
					c.Request().Method,
					c.Path(),
					strconv.Itoa(c.Response().Status),
				).
				Inc()

			monitoring.HTTPRequestDuration.
				WithLabelValues(
					c.Request().Method,
					c.Path(),
				).
				Observe(duration)

			return err
		}
	}
}
