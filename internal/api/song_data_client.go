package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/alserok/music_lib/internal/logger"
	"github.com/alserok/music_lib/internal/service/models"
	"github.com/alserok/music_lib/internal/utils"
	"net/http"
	"net/http/httptest"
	"time"
)

type SongDataAPIClient interface {
	GetSongData(ctx context.Context, group string, song string) (models.SongData, error)
}

func NewSongDataClient(addr string) *songDataClient {
	return &songDataClient{
		addr: addr,
		cl:   http.DefaultClient,
	}
}

const (
	pathInfo = "/info"
)

type songDataClient struct {
	addr string

	cl *http.Client
}

func (s *songDataClient) GetSongData(ctx context.Context, group string, song string) (models.SongData, error) {
	logger.ExtractLogger(ctx).Debug("SongDataAPI sending request", logger.WithArg("id", logger.ExtractIdentifier(ctx)))

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", s.addr, pathInfo), nil)

	query := req.URL.Query()
	query.Set("group", group)
	query.Set("song", song)
	req.URL.RawQuery = query.Encode()

	s.cl.Timeout = 1 * time.Second
	res, err := s.cl.Do(req)
	if err != nil {
		return models.SongData{}, utils.NewError(err.Error(), utils.Internal)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	logger.ExtractLogger(ctx).Debug("SongDataAPI response",
		logger.WithArg("id", logger.ExtractIdentifier(ctx)),
		logger.WithArg("res_status", res.StatusCode),
	)

	switch res.StatusCode {
	case http.StatusOK:
		var songData models.SongData
		if err = json.NewDecoder(res.Body).Decode(&songData); err != nil {
			return models.SongData{}, utils.NewError(err.Error(), utils.Internal)
		}

		return songData, nil
	case http.StatusBadRequest:
		return models.SongData{}, utils.NewError("api request failed", utils.BadRequest)
	default:
		// 5xx, 404 etc. => may be implemented some additional logic
		return models.SongData{}, nil
	}
}
