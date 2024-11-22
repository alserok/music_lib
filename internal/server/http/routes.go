package http

import (
	"github.com/alserok/music_lib/internal/server/http/middleware"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
)

func setupRoutes(s *echo.Echo, h handler) {
	v1 := s.Group("/v1")
	v1.Use(middleware.WithErrorHandler, middleware.WithLogger(h.log))
	v1.GET("/swagger/*", echoSwagger.WrapHandler)

	get := v1.Group("/get")
	get.GET("/songs", h.GetSongs)
	get.GET("/songs/:id", h.GetSongText)

	del := v1.Group("/del")
	del.DELETE("/:id", h.DeleteSong)

	edit := v1.Group("/edit")
	edit.PUT("/", h.EditSong)

	create := v1.Group("/new")
	create.POST("/song", h.CreateSong)
}
