package server

import (
	"github.com/alserok/music_lib/internal/server/http"
	"github.com/alserok/music_lib/internal/service"
)

type Server interface {
	MustServe(port string)
}

const (
	HTTP = iota
)

func New(serverType uint, srvc service.Service) Server {
	switch serverType {
	case HTTP:
		return http.NewServer(srvc)
	default:
		panic("invalid server type")
	}
}
