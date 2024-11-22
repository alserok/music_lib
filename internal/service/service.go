package service

import (
	"context"
	"fmt"
	"github.com/alserok/music_lib/internal/api"
	"github.com/alserok/music_lib/internal/db"
	"github.com/alserok/music_lib/internal/logger"
	"github.com/alserok/music_lib/internal/service/models"
	"github.com/alserok/music_lib/internal/utils"
	"github.com/google/uuid"
	"strings"
)

type Service interface {
	CreateSong(ctx context.Context, song models.Song) error
	EditSong(ctx context.Context, song models.Song) error
	DeleteSong(ctx context.Context, songID string) error
	GetSongText(ctx context.Context, songID string, lim, off int) (string, error)
	GetSongs(ctx context.Context, filter models.SongFilter) ([]models.Song, error)
}

type Clients struct {
	SongDataAPIClient api.SongDataAPIClient
}

func New(repo db.Repository, cls *Clients) *service {
	return &service{
		repo:              repo,
		songDataAPIClient: cls.SongDataAPIClient,
	}
}

type service struct {
	repo db.Repository

	songDataAPIClient api.SongDataAPIClient
}

func (s *service) CreateSong(ctx context.Context, song models.Song) error {
	logger.ExtractLogger(ctx).
		Debug("service received CreateSong",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)
	defer logger.ExtractLogger(ctx).
		Debug("service passed CreateSong",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	songData, err := s.songDataAPIClient.GetSongData(ctx, song.Group, song.Song)
	if err != nil {
		return fmt.Errorf("client failed to get song data: %w", err)
	}

	song.SongID = uuid.NewString()
	song.Data = songData

	if err = s.repo.CreateSong(ctx, song); err != nil {
		return fmt.Errorf("repo failed to create song: %w", err)
	}

	return nil
}

func (s *service) EditSong(ctx context.Context, song models.Song) error {
	logger.ExtractLogger(ctx).
		Debug("service received EditSong",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)
	defer logger.ExtractLogger(ctx).
		Debug("service passed EditSong",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	if err := s.repo.EditSong(ctx, song); err != nil {
		return fmt.Errorf("repo failed to edit song: %w", err)
	}

	return nil
}

func (s *service) DeleteSong(ctx context.Context, songID string) error {
	logger.ExtractLogger(ctx).
		Debug("service received DeleteSong",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)
	defer logger.ExtractLogger(ctx).
		Debug("service passed DeleteSong",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	if err := s.repo.DeleteSong(ctx, songID); err != nil {
		return fmt.Errorf("repo failed to delete song: %w", err)
	}

	return nil
}

func (s *service) GetSongText(ctx context.Context, songID string, lim, off int) (string, error) {
	logger.ExtractLogger(ctx).
		Debug("service received GetSongText",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)
	defer logger.ExtractLogger(ctx).
		Debug("service passed GetSongText",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	text, err := s.repo.GetSongText(ctx, songID)
	if err != nil {
		return "", fmt.Errorf("repo failed to get song text: %w", err)
	}

	couplets := strings.Split(text, "\n\n")

	if off+lim > len(couplets) {
		return "", utils.NewError(
			fmt.Sprintf("invalid pagination parameters: number of couplets: %d last_requested_couplet_index: %d", len(couplets), off+lim-1), utils.BadRequest)
	}

	return strings.Join(couplets[off:lim+off], "\n\n"), nil
}

func (s *service) GetSongs(ctx context.Context, filter models.SongFilter) ([]models.Song, error) {
	logger.ExtractLogger(ctx).
		Debug("service received GetSongs",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)
	defer logger.ExtractLogger(ctx).
		Debug("service passed GetSongs",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	songs, err := s.repo.GetSongs(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("repo failed to get songs: %w", err)
	}

	return songs, nil
}
