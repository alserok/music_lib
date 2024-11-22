package logger

import (
	"io"
	"log/slog"
	"os"
)

func NewSlog(env string) *slogLogger {
	var l *slog.Logger

	switch env {
	case "PROD":
		file, err := os.OpenFile("production.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic("failed to open log file: " + err.Error())
		}
		l = slog.New(slog.NewJSONHandler(io.MultiWriter(file), &slog.HandlerOptions{Level: slog.LevelInfo}))
	case "DEV":
		l = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		l = slog.New(slog.NewJSONHandler(io.MultiWriter(os.Stdout), &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return &slogLogger{l}
}

type slogLogger struct {
	l *slog.Logger
}

func (l *slogLogger) Info(msg string, args ...Arg) {
	l.l.Info(msg, args)
}

func (l *slogLogger) Error(msg string, args ...Arg) {
	l.l.Error(msg, args)
}

func (l *slogLogger) Debug(msg string, args ...Arg) {
	l.l.Debug(msg, args)
}
