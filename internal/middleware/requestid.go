package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RequestID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := uuid.New().String()
		c.Response().Header().Set("X-Request-ID", id)
		return next(c)
	}
}
