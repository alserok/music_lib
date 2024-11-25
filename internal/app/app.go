package app

import (
	"github.com/alserok/music_lib/internal/api"
	"github.com/alserok/music_lib/internal/config"
	"github.com/alserok/music_lib/internal/db/postgres"
	"github.com/alserok/music_lib/internal/logger"
	"github.com/alserok/music_lib/internal/server"
	"github.com/alserok/music_lib/internal/service"

	_ "github.com/alserok/music_lib/docs"
)

func MustStart(cfg *config.Config) {
	log := logger.NewSlog(cfg.Env)
	log.Info("starting server")
	defer log.Info("server was stopped")

	conn := postgres.MustConnect(cfg.DB.DSN())
	defer func() {
		_ = conn.Close()
	}()

	songDataClient := api.NewSongDataClient(cfg.Clients.SongDataAPIAddr)

	repo := postgres.NewRepository(conn)
	srvc := service.New(repo, &service.Clients{SongDataAPIClient: songDataClient})
	srvr := server.New(server.HTTP, srvc, log)

	log.Info("server is running", logger.WithArg("port", cfg.Port))
	srvr.MustServe(cfg.Port)
}
