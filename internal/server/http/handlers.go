package http

import (
	"fmt"
	"github.com/alserok/music_lib/internal/logger"
	"github.com/alserok/music_lib/internal/service"
	"github.com/alserok/music_lib/internal/service/models"
	"github.com/alserok/music_lib/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type handler struct {
	srvc service.Service
	log  logger.Logger
}

func newHandler(srvc service.Service, log logger.Logger) handler {
	return handler{srvc: srvc, log: log}
}

// @Summary GetSongs
// @Description Get a list of songs with optional filters
// @Tags songs
// @Accept json
// @Produce json
// @Param limit query int true "Limit of songs to return"
// @Param offset query int true "Offset for pagination"
// @Param songID query string false "Filter by songID"
// @Param group query string false "Filter by group"
// @Param song query string false "Filter by song name"
// @Param text query string false "Filter by text"
// @Param releaseDate query string false "Filter by release date"
// @Param link query string false "Filter by link"
// @Success 200 {array} models.Song "Success"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal error"
// @Router /get/songs [get]
func (h *handler) GetSongs(c echo.Context) error {
	logger.ExtractLogger(c.Request().Context()).
		Debug("received GetSongs request",
			logger.WithArg("id", logger.ExtractIdentifier(c.Request().Context())),
		)

	lim, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return utils.NewError("failed to parse limit", utils.BadRequest)
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return utils.NewError("failed to parse offset", utils.BadRequest)
	}

	filter := models.SongFilter{
		Lim:         lim,
		Off:         offset,
		SongID:      c.QueryParam("songID"),
		Group:       c.QueryParam("group"),
		Song:        c.QueryParam("song"),
		Text:        c.QueryParam("text"),
		ReleaseDate: c.QueryParam("releaseDate"),
		Link:        c.QueryParam("link"),
	}

	songs, err := h.srvc.GetSongs(c.Request().Context(), filter)
	if err != nil {
		return fmt.Errorf("failed to get songs: %w", err)
	}

	logger.ExtractLogger(c.Request().Context()).
		Debug("passed GetSongs request",
			logger.WithArg("id", logger.ExtractIdentifier(c.Request().Context())),
		)

	return c.JSON(http.StatusOK, map[string]interface{}{"songs": songs})
}

// @Summary GetSongText
// @Description Get the text of a specific song
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Param limit query int true "Limit of text entries to return"
// @Param offset query int true "Offset for pagination"
// @Success 200 {string} string "Success"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Not found"
// @Failure 500 {object} string "Internal error"
// @Router /get/songs/{id} [get]
func (h *handler) GetSongText(c echo.Context) error {
	logger.ExtractLogger(c.Request().Context()).
		Debug("received GetSongText request",
			logger.WithArg("id", logger.ExtractIdentifier(c.Request().Context())),
		)

	lim, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return utils.NewError("failed to parse limit", utils.BadRequest)
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return utils.NewError("failed to parse offset", utils.BadRequest)
	}

	songID := c.Param("id")
	if songID == "" {
		return utils.NewError("songID is required", utils.BadRequest)
	}

	text, err := h.srvc.GetSongText(c.Request().Context(), songID, lim, offset)
	if err != nil {
		return fmt.Errorf("failed to get song text: %w", err)
	}

	logger.ExtractLogger(c.Request().Context()).
		Debug("passed GetSongText request",
			logger.WithArg("id", logger.ExtractIdentifier(c.Request().Context())),
		)

	return c.JSON(http.StatusOK, map[string]interface{}{"text": text})
}

// @Summary DeleteSong
// @Description Delete a specific song
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Success 200 {object} interface{} "Success"
// @Failure 400 {object} string "Bad request"
// @Failure 404 {object} string "Not found"
// @Failure 500 {object} string "Internal error"
// @Router /del/{id} [delete]
func (h *handler) DeleteSong(c echo.Context) error {
	logger.ExtractLogger(c.Request().Context()).
		Debug("received DeleteSong request",
			logger.WithArg("id", logger.ExtractIdentifier(c.Request().Context())),
		)

	songID := c.Param("id")

	err := h.srvc.DeleteSong(c.Request().Context(), songID)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}

	logger.ExtractLogger(c.Request().Context()).
		Debug("passed DeleteSong request",
			logger.WithArg("id", logger.ExtractIdentifier(c.Request().Context())),
		)

	return c.JSON(http.StatusOK, nil)
}

// @Summary EditSong
// @Description Edit a specific song
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song details"
// @Success 200 {object} interface{} "Success"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal error"
// @Router /edit/ [put]
func (h *handler) EditSong(c echo.Context) error {
	logger.ExtractLogger(c.Request().Context()).
		Debug("received EditSong request",
			logger.WithArg("id", logger.ExtractIdentifier(c.Request().Context())),
		)

	var song models.Song
	if err := c.Bind(&song); err != nil {
		return utils.NewError(err.Error(), utils.BadRequest)
	}

	if err := h.srvc.EditSong(c.Request().Context(), song); err != nil {
		return fmt.Errorf("failed to create song: %w", err)
	}

	logger.ExtractLogger(c.Request().Context()).
		Debug("passed EditSong request",
			logger.WithArg("id", logger.ExtractIdentifier(c.Request().Context())),
		)

	return c.JSON(http.StatusOK, nil)
}

// @Summary CreateSong
// @Description Add a new song to the library
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.NewSong true "Song details"
// @Success 201 {object} interface{} "Created"
// @Failure 400 {object} string "Bad request"
// @Failure 500 {object} string "Internal error"
// @Router /new/song [post]
func (h *handler) CreateSong(c echo.Context) error {
	logger.ExtractLogger(c.Request().Context()).
		Debug("received CreateSong request",
			logger.WithArg("id", logger.ExtractIdentifier(c.Request().Context())),
		)

	var song models.NewSong
	if err := c.Bind(&song); err != nil {
		return utils.NewError(err.Error(), utils.BadRequest)
	}

	if err := h.srvc.CreateSong(c.Request().Context(), models.Song{Song: song.Song, Group: song.Group}); err != nil {
		return fmt.Errorf("failed to create song: %w", err)
	}

	logger.ExtractLogger(c.Request().Context()).
		Debug("passed CreateSong request",
			logger.WithArg("id", logger.ExtractIdentifier(c.Request().Context())),
		)

	return c.JSON(http.StatusCreated, nil)
}
