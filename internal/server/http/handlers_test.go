package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/alserok/music_lib/internal/logger"
	"github.com/alserok/music_lib/internal/mocks"
	"github.com/alserok/music_lib/internal/service"
	"github.com/alserok/music_lib/internal/service/models"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestHandlersSuite(t *testing.T) {
	suite.Run(t, new(HTTPHandlersSuite))
}

type HTTPHandlersSuite struct {
	suite.Suite
	ctrl *gomock.Controller

	e *echo.Echo

	handler handler
	logger  *mocks.MockLogger
	repo    *mocks.MockRepository
	api     *mocks.MockSongDataAPIClient
}

func (suite *HTTPHandlersSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.repo = mocks.NewMockRepository(suite.ctrl)
	suite.api = mocks.NewMockSongDataAPIClient(suite.ctrl)
	suite.logger = mocks.NewMockLogger(suite.ctrl)
	suite.e = echo.New()

	suite.handler = handler{
		srvc: service.New(suite.repo, &service.Clients{SongDataAPIClient: suite.api}),
		log:  suite.logger,
	}
}

func (suite *HTTPHandlersSuite) TeardownTest() {
	suite.ctrl.Finish()
	suite.Require().NoError(suite.e.Shutdown(context.Background()))
}

func (suite *HTTPHandlersSuite) TestGetSongs() {
	filter := models.SongFilter{
		SongID:      "id",
		Group:       "group",
		Song:        "song",
		ReleaseDate: "11.11.2011",
		Text:        "text",
		Link:        "link",
		Lim:         1,
		Off:         1,
	}
	songs := []models.Song{
		{
			Song:   filter.Song,
			Group:  filter.Group,
			SongID: filter.SongID,
			Data: models.SongData{
				ReleaseDate: filter.ReleaseDate,
				Text:        filter.Text,
				Link:        filter.Link,
			},
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(logger.WrapLogger(req.Context(), suite.logger))
	req = req.WithContext(logger.WrapIdentifier(req.Context()))
	query := req.URL.Query()
	query.Set("limit", strconv.Itoa(filter.Lim))
	query.Set("offset", strconv.Itoa(filter.Off))
	query.Set("songID", filter.SongID)
	query.Set("group", filter.Group)
	query.Set("song", filter.Song)
	query.Set("releaseDate", filter.ReleaseDate)
	query.Set("text", filter.Text)
	query.Set("link", filter.Link)
	req.URL.RawQuery = query.Encode()
	rec := httptest.NewRecorder()

	suite.repo.EXPECT().
		GetSongs(gomock.Any(), gomock.Eq(filter)).
		Return(songs, nil).
		Times(1)

	suite.logger.EXPECT().
		Debug(gomock.Any(), gomock.Eq(logger.Arg{Key: "id", Val: logger.ExtractIdentifier(req.Context())})).
		AnyTimes()

	c := suite.e.NewContext(req, rec)
	suite.Require().NoError(suite.handler.GetSongs(c))
	suite.Equal(http.StatusOK, rec.Code)

	var res map[string][]models.Song
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &res))
	suite.Equal(res["songs"], songs)
}

func (suite *HTTPHandlersSuite) TestCreateSong() {
	song := models.Song{
		Song:   "song",
		Group:  "group",
		SongID: "id",
	}

	songData := models.SongData{
		ReleaseDate: "11.11.2011",
		Text:        "text",
		Link:        "link",
	}

	b, err := json.Marshal(song)
	suite.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req = req.WithContext(logger.WrapLogger(req.Context(), suite.logger))
	req = req.WithContext(logger.WrapIdentifier(req.Context()))
	rec := httptest.NewRecorder()

	suite.repo.EXPECT().
		CreateSong(gomock.Any(), gomock.Any()).
		Return(nil).
		Times(1)

	suite.logger.EXPECT().
		Debug(gomock.Any(), gomock.Eq(logger.Arg{Key: "id", Val: logger.ExtractIdentifier(req.Context())})).
		AnyTimes()

	suite.api.EXPECT().
		GetSongData(gomock.Any(), gomock.Eq(song.Group), gomock.Eq(song.Song)).
		Return(songData, nil).
		Times(1)

	c := suite.e.NewContext(req, rec)
	suite.Require().NoError(suite.handler.CreateSong(c))
	suite.Equal(http.StatusCreated, rec.Code)
}

func (suite *HTTPHandlersSuite) TestEditSong() {
	song := models.Song{
		Song:   "song",
		Group:  "group",
		SongID: "id",
		Data: models.SongData{
			ReleaseDate: "11.11.2011",
			Text:        "text",
			Link:        "link",
		},
	}

	b, err := json.Marshal(song)
	suite.Require().NoError(err)

	req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(b))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req = req.WithContext(logger.WrapLogger(req.Context(), suite.logger))
	req = req.WithContext(logger.WrapIdentifier(req.Context()))
	rec := httptest.NewRecorder()

	suite.repo.EXPECT().
		EditSong(gomock.Any(), gomock.Eq(song)).
		Return(nil).
		Times(1)

	suite.logger.EXPECT().
		Debug(gomock.Any(), gomock.Eq(logger.Arg{Key: "id", Val: logger.ExtractIdentifier(req.Context())})).
		AnyTimes()

	c := suite.e.NewContext(req, rec)
	suite.Require().NoError(suite.handler.EditSong(c))
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *HTTPHandlersSuite) TestDeleteSong() {
	songID := "id"

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(logger.WrapLogger(req.Context(), suite.logger))
	req = req.WithContext(logger.WrapIdentifier(req.Context()))
	rec := httptest.NewRecorder()

	suite.repo.EXPECT().
		DeleteSong(gomock.Any(), gomock.Eq(songID)).
		Return(nil).
		Times(1)

	suite.logger.EXPECT().
		Debug(gomock.Any(), gomock.Eq(logger.Arg{Key: "id", Val: logger.ExtractIdentifier(req.Context())})).
		AnyTimes()

	c := suite.e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(songID)
	suite.Require().NoError(suite.handler.DeleteSong(c))
	suite.Equal(http.StatusOK, rec.Code)
}

func (suite *HTTPHandlersSuite) TestGetSongText() {
	couplets := []string{"c0", "c1", "c2", "c3"}
	song := models.Song{
		Song:   "song",
		Group:  "group",
		SongID: "id",
		Data: models.SongData{
			ReleaseDate: "11.11.2011",
			Text:        strings.Join(couplets, "\n\n"),
			Link:        "link",
		},
	}
	lim, off := 1, 2

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(logger.WrapLogger(req.Context(), suite.logger))
	req = req.WithContext(logger.WrapIdentifier(req.Context()))
	query := req.URL.Query()
	query.Set("limit", strconv.Itoa(lim))
	query.Set("offset", strconv.Itoa(off))
	req.URL.RawQuery = query.Encode()
	rec := httptest.NewRecorder()

	suite.repo.EXPECT().
		GetSongText(gomock.Any(), gomock.Eq(song.SongID), gomock.Eq(lim), gomock.Eq(off)).
		Return(song.Data.Text, nil).
		Times(1)

	suite.logger.EXPECT().
		Debug(gomock.Any(), gomock.Eq(logger.Arg{Key: "id", Val: logger.ExtractIdentifier(req.Context())})).
		AnyTimes()

	c := suite.e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(song.SongID)
	suite.Require().NoError(suite.handler.GetSongText(c))
	suite.Equal(http.StatusOK, rec.Code)

	var res map[string]string
	suite.NoError(json.Unmarshal(rec.Body.Bytes(), &res))
	suite.Equal(strings.Join(couplets[off:off+lim], "\n\n"), res["text"])
}
