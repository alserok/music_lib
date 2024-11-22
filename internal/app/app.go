package app

import (
	"github.com/alserok/music_lib/internal/config"
	"github.com/alserok/music_lib/internal/db/postgres"
	"github.com/alserok/music_lib/internal/logger"
	"github.com/alserok/music_lib/internal/server"
	"github.com/alserok/music_lib/internal/service"
)

func MustStart(cfg *config.Config) {
	log := logger.NewSlog(cfg.Env)
	log.Info("starting server")
	defer log.Info("server was stopped")

	conn := postgres.MustConnect(cfg.DB.DSN())
	defer func() {
		_ = conn.Close()
	}()

	repo := postgres.NewRepository(conn)
	srvc := service.New(repo)
	srvr := server.New(server.HTTP, srvc)

	srvr.MustServe(cfg.Port)
}
