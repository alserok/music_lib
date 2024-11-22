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
	arguments := make([]any, 0, len(args))
	for _, arg := range args {
		arguments = append(arguments, slog.Attr{
			Key:   arg.Key,
			Value: slog.AnyValue(arg.Val),
		})
	}

	l.l.Info(msg, arguments...)
}

func (l *slogLogger) Error(msg string, args ...Arg) {
	arguments := make([]any, 0, len(args))
	for _, arg := range args {
		arguments = append(arguments, slog.Attr{
			Key:   arg.Key,
			Value: slog.AnyValue(arg.Val),
		})
	}

	l.l.Error(msg, arguments...)
}

func (l *slogLogger) Debug(msg string, args ...Arg) {
	arguments := make([]any, 0, len(args))
	for _, arg := range args {
		arguments = append(arguments, slog.Attr{
			Key:   arg.Key,
			Value: slog.AnyValue(arg.Val),
		})
	}

	l.l.Debug(msg, arguments...)
}
