package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/alserok/music_lib/internal/logger"
	"github.com/alserok/music_lib/internal/service/models"
	"github.com/alserok/music_lib/internal/utils"
	"github.com/jmoiron/sqlx"
	"time"
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

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	q := `INSERT INTO songs (id, song, release_date, text, link) VALUES ($1, $2, $3, $4, $5)`

	_, err = tx.ExecContext(ctx, q, song.SongID, song.Song, song.Data.ReleaseDate, song.Data.Text, song.Data.Link)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}

	q = `INSERT INTO group_songs (song_id, group_name) VALUES ($1, $2)`

	_, err = tx.ExecContext(ctx, q, song.SongID, song.Group)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}

	if err = tx.Commit(); err != nil {
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

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	q := `UPDATE songs SET song = $2, release_date = $3, text = $4, link = $5 WHERE id = $1`

	_, err = tx.ExecContext(ctx, q, song.SongID, song.Song, song.Data.ReleaseDate, song.Data.Text, song.Data.Link)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}

	q = `UPDATE group_songs SET group_name = $2 WHERE song_id = $1`

	_, err = tx.ExecContext(ctx, q, song.SongID, song.Group)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}

	if err = tx.Commit(); err != nil {
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

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	q := `DELETE FROM group_songs WHERE song_id = $1`

	_, err = r.db.ExecContext(ctx, q, songID)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}

	q = `DELETE FROM songs WHERE id = $1`

	_, err = tx.ExecContext(ctx, q, songID)
	if err != nil {
		return utils.NewError(err.Error(), utils.Internal)
	}

	if err = tx.Commit(); err != nil {
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

	q := `SELECT text FROM songs WHERE id = $1 LIMIT 1`

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

	q := `SELECT 
				group_songs.song_id as id, 
				group_songs.group_name, 
				songs.song, 
				songs.release_date, 
				songs.text, 
				songs.link 
			FROM songs INNER JOIN group_songs ON songs.id = group_songs.song_id
      WHERE 
          (group_songs.song_id = $1 OR $1 = '') AND
          (group_songs.group_name LIKE '%' || $2 || '%' OR $2 = '') AND
          (songs.song LIKE '%' || $3 || '%' OR $3 = '') AND
          (songs.release_date = $4 OR $4 IS NULL) AND
          (songs.text LIKE '%' || $5 || '%' OR $5 = '') AND
          (songs.link = $6 OR $6 = '')
      OFFSET $7 LIMIT $8`

	rows, err := r.db.QueryxContext(ctx, q, filter.SongID, filter.Group, filter.Song, filter.ReleaseDate, filter.Text, filter.Link, filter.Off, filter.Lim)
	if err != nil {
		return nil, utils.NewError(err.Error(), utils.Internal)
	}

	songs := make([]models.Song, 0, filter.Lim)
	for rows.Next() {
		var fullSongData struct {
			SongID      string    `json:"songID" db:"id"`
			Group       string    `json:"group" db:"group_name"`
			Song        string    `json:"song"`
			ReleaseDate time.Time `json:"releaseDate" db:"release_date"`
			Text        string    `json:"text"`
			Link        string    `json:"link"`
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
