package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alserok/music_lib/internal/logger"
	"github.com/alserok/music_lib/internal/service/models"
	"github.com/alserok/music_lib/internal/utils"
	"github.com/jmoiron/sqlx"
)

func NewRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) CreateSong(ctx context.Context, song models.Song) error {
	logger.ExtractLogger(ctx).
		Debug("repo received GetSongText",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	q := `INSERT INTO songs (song_id, group_name, song, release_date, text, link) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.ExecContext(ctx, q, song.SongID, song.Group, song.Song, song.Data.ReleaseDate, song.Data.Text, song.Data.Link)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}

	logger.ExtractLogger(ctx).
		Debug("repo passed CreateSong",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	return nil
}

func (r *repository) EditSong(ctx context.Context, song models.Song) error {
	logger.ExtractLogger(ctx).
		Debug("repo received EditSong",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	q := `UPDATE songs SET group_name = $2, song = $3, release_date = $4, text = $5, link = $6 WHERE song_id = $1`

	_, err := r.db.ExecContext(ctx, q, song.SongID, song.Group, song.Song, song.Data.ReleaseDate, song.Data.Text, song.Data.Link)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}

	logger.ExtractLogger(ctx).
		Debug("repo passed EditSong",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	return nil
}

func (r *repository) DeleteSong(ctx context.Context, songID string) error {
	logger.ExtractLogger(ctx).
		Debug("repo received DeleteSong",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	q := `DELETE FROM songs WHERE song_id = $1`

	_, err := r.db.ExecContext(ctx, q, songID)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}

	logger.ExtractLogger(ctx).
		Debug("repo passed DeleteSong",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	return nil
}

func (r *repository) GetSongText(ctx context.Context, songID string) (string, error) {
	logger.ExtractLogger(ctx).
		Debug("repo received GetSongText",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	q := `SELECT text FROM songs WHERE song_id = $1 LIMIT 1`

	var text string
	if err := r.db.QueryRowxContext(ctx, q, songID).Scan(&text); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", utils.NewError("song not found", utils.NotFound)
		}
		return "", utils.NewError(err.Error(), utils.Internal)
	}

	logger.ExtractLogger(ctx).
		Debug("repo passed GetSongText",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	return text, nil
}

func (r *repository) GetSongs(ctx context.Context, filter models.SongFilter) ([]models.Song, error) {
	logger.ExtractLogger(ctx).
		Debug("repo received GetSongs",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	q := `SELECT * FROM songs 
      WHERE 
          (song_id = $1 OR $1 = '') AND
          (group_name LIKE '%' || $2 || '%' OR $2 = '') AND
          (song LIKE '%' || $3 || '%' OR $3 = '') AND
          (release_date = $4 OR $4 = '') AND
          (text LIKE '%' || $5 || '%' OR $5 = '') AND
          (link = $6 OR $6 = '')
      OFFSET $8 LIMIT $7`

	rows, err := r.db.QueryxContext(ctx, q,
		filter.SongID, filter.Group, filter.Song, filter.ReleaseDate, filter.Text, filter.Link, filter.Lim, filter.Off)
	if err != nil {
		return nil, utils.NewError(err.Error(), utils.Internal)
	}

	songs := make([]models.Song, 0, filter.Lim)
	for rows.Next() {
		var fullSongData struct {
			SongID      string `json:"songID" db:"song_id"`
			Group       string `json:"group" db:"group_name"`
			Song        string `json:"song"`
			ReleaseDate string `json:"releaseDate" db:"release_date"`
			Text        string `json:"text"`
			Link        string `json:"link"`
		}
		if err = rows.StructScan(&fullSongData); err != nil {
			// may not return an error and continue with the other songs
			return nil, utils.NewError(err.Error(), utils.Internal)
		}

		songs = append(songs, models.Song{
			SongID: fullSongData.SongID,
			Group:  fullSongData.Group,
			Song:   fullSongData.Song,
			Data: models.SongData{
				ReleaseDate: fullSongData.ReleaseDate,
				Text:        fullSongData.Text,
				Link:        fullSongData.Link,
			},
		})
	}

	logger.ExtractLogger(ctx).
		Debug("repo passed GetSongs",
			logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		)

	return songs, nil
}
