package service

import "github.com/alserok/music_lib/internal/db"

type Service interface {
}

func New(repo db.Repository) *service {
	return &service{
		repo: repo,
	}
}

type service struct {
	repo db.Repository
}
