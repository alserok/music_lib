package middleware

import (
	"github.com/alserok/music_lib/internal/utils"
	"github.com/labstack/echo/v4"
)

func WithErrorHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err != nil {
			code, msg := utils.FromErrorToHTTP(c.Request().Context(), err)
			_ = c.JSON(code, map[string]interface{}{
				"error": msg,
			})
		}

		return nil
	}
}
