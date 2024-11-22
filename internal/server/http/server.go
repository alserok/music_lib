package http

import (
	"context"
	"github.com/alserok/music_lib/internal/service"
	"github.com/labstack/echo/v4"
	"os/signal"
	"syscall"
)

func NewServer(srvc service.Service) *server {
	return &server{
		srvc: srvc,
		serv: echo.New(),
	}
}

type server struct {
	srvc service.Service

	serv *echo.Echo
}

func (s server) MustServe(port string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	setupRoutes(s.serv)

	go func() {
		if err := s.serv.Start(port); err != nil {
			panic("failed to start server: " + err.Error())
		}
	}()

	<-ctx.Done()
}
