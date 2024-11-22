package db

import (
	"context"
	"github.com/alserok/music_lib/internal/service/models"
)

type Repository interface {
	CreateSong(ctx context.Context, song models.Song) error
	EditSong(ctx context.Context, song models.Song) error
	DeleteSong(ctx context.Context, songID string) error
	GetSongText(ctx context.Context, songID string, lim, off int) (string, error)
	GetSongs(ctx context.Context, filter models.SongFilter) ([]models.Song, error)
}
