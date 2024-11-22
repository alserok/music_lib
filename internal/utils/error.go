package utils

import (
	"context"
	"errors"
	"github.com/alserok/music_lib/internal/logger"
	"net/http"
)

type err struct {
	msg  string
	code int
}

func (e *err) Error() string {
	return e.msg
}

const (
	Internal = iota
	BadRequest
)

func NewError(msg string, code int) error {
	return &err{
		msg:  msg,
		code: code,
	}
}

func FromErrorToHTTP(ctx context.Context, in error) (int, string) {
	l := logger.ExtractLogger(ctx)

	var e *err
	if errors.As(in, &e) {
		l.Error("unknown error", logger.WithArg("error", in.Error()))
		return http.StatusInternalServerError, "internal server error"
	}

	switch e.code {
	case Internal:
		l.Error("internal error", logger.WithArg("error", in.Error()))
		return http.StatusInternalServerError, "internal server error"
	case BadRequest:
		return http.StatusBadRequest, e.msg
	default:
		l.Error("unknown error code", logger.WithArg("code", e.code))
		return http.StatusInternalServerError, "internal server error"
	}
}
