package middleware

import (
	"github.com/alserok/music_lib/internal/logger"
	"github.com/labstack/echo/v4"
)

func WithRecovery(log logger.Logger) func(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(handlerFuncFunc echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovery", logger.WithArg("recover", rec))
				}
			}()

			return handlerFuncFunc(c)
		}
	}
}
