package middleware

import (
	"github.com/alserok/music_lib/internal/logger"
	"github.com/labstack/echo/v4"
)

func WithLogger(log logger.Logger) func(echo.HandlerFunc) echo.HandlerFunc {
	return func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = logger.WrapLogger(ctx, log)
			ctx = logger.WrapIdentifier(ctx)
			c.SetRequest(c.Request().WithContext(ctx))

			return handlerFunc(c)
		}
	}
}
