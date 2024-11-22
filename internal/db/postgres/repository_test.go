package postgres

import (
	"context"
	"github.com/alserok/music_lib/internal/config"
	"github.com/alserok/music_lib/internal/logger"
	"github.com/alserok/music_lib/internal/mocks"
	"github.com/alserok/music_lib/internal/service/models"
	"github.com/docker/go-connections/nat"
	"github.com/golang/mock/gomock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
	"time"
)

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

type RepositorySuite struct {
	suite.Suite

	ctrl   *gomock.Controller
	logger *mocks.MockLogger

	repo      *repository
	conn      *sqlx.DB
	container *postgres.PostgresContainer
}

func (suite *RepositorySuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.logger = mocks.NewMockLogger(suite.ctrl)

	suite.conn, suite.container = newPostgresDB(&suite.Suite)
	suite.repo = NewRepository(suite.conn)
}

func (suite *RepositorySuite) TeardownTest() {
	suite.ctrl.Finish()
	suite.Require().NoError(suite.conn.Close())
	suite.Require().NoError(suite.container.Terminate(context.Background()))
}

func (suite *RepositorySuite) TestGetSongs() {
	songs := []models.Song{
		{
			SongID: "id1",
			Song:   "song1",
			Group:  "group1",
			Data: models.SongData{
				ReleaseDate: "11.11.2011",
				Text:        "song text 1",
				Link:        "link1",
			},
		},
		{
			SongID: "id2",
			Song:   "song2",
			Group:  "group2",
			Data: models.SongData{
				ReleaseDate: "11.11.2011",
				Text:        "song text 2",
				Link:        "link2",
			},
		},
	}
	for _, song := range songs {
		_, err := suite.conn.Exec(`INSERT INTO songs (song_id, group_name, song, release_date, text, link) VALUES ($1, $2, $3, $4, $5,$6)`,
			song.SongID, song.Group, song.Song, song.Data.ReleaseDate, song.Data.Text, song.Data.Link,
		)
		suite.Require().NoError(err)
	}

	suite.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	ctx := logger.WrapLogger(context.Background(), suite.logger)
	ctx = logger.WrapIdentifier(ctx)

	// get all songs
	res, err := suite.repo.GetSongs(ctx, models.SongFilter{Lim: len(songs)})
	suite.Require().NoError(err)
	suite.Require().Len(res, 2)
	for i, song := range res {
		suite.Require().Equal(song, res[i])
	}

	// get 1 song
	res, err = suite.repo.GetSongs(ctx, models.SongFilter{Song: "1", Lim: 1})
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)
	suite.Require().Equal(songs[0].Song, res[0].Song)

	res, err = suite.repo.GetSongs(ctx, models.SongFilter{Group: "1", Lim: 1})
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)
	suite.Require().Equal(songs[0].Song, res[0].Song)

	res, err = suite.repo.GetSongs(ctx, models.SongFilter{Text: "1", Lim: 1})
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)
	suite.Require().Equal(songs[0].Song, res[0].Song)

	res, err = suite.repo.GetSongs(ctx, models.SongFilter{SongID: "id1", Lim: 1})
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)
	suite.Require().Equal(songs[0].Song, res[0].Song)

	res, err = suite.repo.GetSongs(ctx, models.SongFilter{Lim: 1})
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)
	suite.Require().Equal(songs[0].Song, res[0].Song)

	// get 2 song
	res, err = suite.repo.GetSongs(ctx, models.SongFilter{Song: "2", Lim: 1})
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)
	suite.Require().Equal(songs[1].Song, res[0].Song)

	res, err = suite.repo.GetSongs(ctx, models.SongFilter{Group: "2", Lim: 1})
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)
	suite.Require().Equal(songs[1].Song, res[0].Song)

	res, err = suite.repo.GetSongs(ctx, models.SongFilter{Text: "2", Lim: 1})
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)
	suite.Require().Equal(songs[1].Song, res[0].Song)

	res, err = suite.repo.GetSongs(ctx, models.SongFilter{SongID: "id2", Lim: 1})
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)
	suite.Require().Equal(songs[1].Song, res[0].Song)

	res, err = suite.repo.GetSongs(ctx, models.SongFilter{Lim: 1, Off: 1})
	suite.Require().NoError(err)
	suite.Require().Len(res, 1)
	suite.Require().Equal(songs[1].Song, res[0].Song)
}

func (suite *RepositorySuite) TestGetSongText() {
	song := models.Song{
		SongID: "id1",
		Song:   "song1",
		Group:  "group1",
		Data: models.SongData{
			ReleaseDate: "11.11.2011",
			Text:        "song text 1\n\nsong text 2",
			Link:        "link1",
		},
	}
	_, err := suite.conn.Exec(`INSERT INTO songs (song_id, group_name, song, release_date, text, link) VALUES ($1, $2, $3, $4, $5,$6)`,
		song.SongID, song.Group, song.Song, song.Data.ReleaseDate, song.Data.Text, song.Data.Link,
	)
	suite.Require().NoError(err)

	suite.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	ctx := logger.WrapLogger(context.Background(), suite.logger)
	ctx = logger.WrapIdentifier(ctx)

	// get all songs
	res, err := suite.repo.GetSongText(ctx, song.SongID)
	suite.Require().NoError(err)
	suite.Require().Equal(song.Data.Text, res)
}

func (suite *RepositorySuite) TestCreateSong() {
	song := models.Song{
		SongID: "id1",
		Song:   "song1",
		Group:  "group1",
		Data: models.SongData{
			ReleaseDate: "11.11.2011",
			Text:        "song text 1\n\nsong text 2",
			Link:        "link1",
		},
	}

	suite.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	ctx := logger.WrapLogger(context.Background(), suite.logger)
	ctx = logger.WrapIdentifier(ctx)

	err := suite.repo.CreateSong(ctx, song)
	suite.Require().NoError(err)

	var res struct {
		SongID      string `json:"songID" db:"song_id"`
		Group       string `json:"group" db:"group_name"`
		Song        string `json:"song"`
		ReleaseDate string `json:"releaseDate" db:"release_date"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	suite.Require().NoError(suite.conn.QueryRowx(`SELECT * FROM songs LIMIT 1`).StructScan(&res))
	suite.Require().Equal(song, models.Song{SongID: res.SongID, Group: res.Group, Song: res.Song,
		Data: models.SongData{
			ReleaseDate: res.ReleaseDate,
			Text:        res.Text,
			Link:        res.Link,
		}})
}

func (suite *RepositorySuite) TestEditSong() {
	song := models.Song{
		SongID: "id1",
		Song:   "song1",
		Group:  "group1",
		Data: models.SongData{
			ReleaseDate: "11.11.2011",
			Text:        "song text 1\n\nsong text 2",
			Link:        "link1",
		},
	}
	_, err := suite.conn.Exec(`INSERT INTO songs (song_id, group_name, song, release_date, text, link) VALUES ($1, $2, $3, $4, $5,$6)`,
		song.SongID, song.Group, song.Song, song.Data.ReleaseDate, song.Data.Text, song.Data.Link,
	)
	suite.Require().NoError(err)

	suite.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	ctx := logger.WrapLogger(context.Background(), suite.logger)
	ctx = logger.WrapIdentifier(ctx)

	song.Song = "edited song title"
	song.Group = "edited group title"
	song.Data.ReleaseDate = "11.11.2012"
	song.Data.Text = "edited song text"
	song.Data.Link = "edited song link"
	err = suite.repo.EditSong(ctx, song)
	suite.Require().NoError(err)

	var res struct {
		SongID      string `json:"songID" db:"song_id"`
		Group       string `json:"group" db:"group_name"`
		Song        string `json:"song"`
		ReleaseDate string `json:"releaseDate" db:"release_date"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	suite.Require().NoError(suite.conn.QueryRowx(`SELECT * FROM songs LIMIT 1`).StructScan(&res))
	suite.Require().Equal(song, models.Song{SongID: res.SongID, Group: res.Group, Song: res.Song,
		Data: models.SongData{
			ReleaseDate: res.ReleaseDate,
			Text:        res.Text,
			Link:        res.Link,
		}})
}

func (suite *RepositorySuite) TestDeleteSong() {
	song := models.Song{
		SongID: "id1",
		Song:   "song1",
		Group:  "group1",
		Data: models.SongData{
			ReleaseDate: "11.11.2011",
			Text:        "song text 1\n\nsong text 2",
			Link:        "link1",
		},
	}
	_, err := suite.conn.Exec(`INSERT INTO songs (song_id, group_name, song, release_date, text, link) VALUES ($1, $2, $3, $4, $5,$6)`,
		song.SongID, song.Group, song.Song, song.Data.ReleaseDate, song.Data.Text, song.Data.Link,
	)
	suite.Require().NoError(err)

	suite.logger.EXPECT().
		Debug(gomock.Any(), gomock.Any()).
		AnyTimes()

	ctx := logger.WrapLogger(context.Background(), suite.logger)
	ctx = logger.WrapIdentifier(ctx)

	err = suite.repo.DeleteSong(ctx, song.SongID)
	suite.Require().NoError(err)

	var res int64
	suite.Require().NoError(suite.conn.QueryRowx(`SELECT count(*) FROM songs`).Scan(&res))
	suite.Require().Equal(int64(0), res)
}

func newPostgresDB(s *suite.Suite) (*sqlx.DB, *postgres.PostgresContainer) {
	ctx := context.Background()
	cfg := config.Postgres{
		Name: "postgres",
		Port: "5432",
		User: "postgres",
		Pass: "postgres",
		Host: "localhost",
	}

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(cfg.Name),
		postgres.WithUsername(cfg.User),
		postgres.WithPassword(cfg.Pass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	s.Require().NoError(err)
	s.Require().NotNil(postgresContainer)
	s.Require().True(postgresContainer.IsRunning())

	port, err := postgresContainer.MappedPort(ctx, nat.Port(cfg.Port+"/tcp"))
	s.Require().NoError(err)
	cfg.Port = port.Port()

	conn := MustConnect(cfg.DSN(), "../migrations")

	return conn, postgresContainer
}
