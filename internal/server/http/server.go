package http

import (
	"context"
	"github.com/alserok/music_lib/internal/logger"
	"github.com/alserok/music_lib/internal/service"
	"github.com/labstack/echo/v4"
	"os/signal"
	"syscall"
)

func NewServer(srvc service.Service, log logger.Logger) *server {
	return &server{
		srvc: srvc,
		serv: echo.New(),
		log:  log,
	}
}

type server struct {
	srvc service.Service

	log logger.Logger

	serv *echo.Echo
}

func (s server) MustServe(port string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	setupRoutes(s.serv, newHandler(s.srvc, s.log))

	go func() {
		if err := s.serv.Start(port); err != nil {
			panic("failed to start server: " + err.Error())
		}
	}()

	<-ctx.Done()
}
