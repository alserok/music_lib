package logger

import (
	"context"
	"github.com/google/uuid"
)

type Logger interface {
	Info(msg string, args ...Arg)
	Error(msg string, args ...Arg)
	Debug(msg string, args ...Arg)
}

func WithArg(key string, value interface{}) Arg {
	return Arg{key, value}
}

type Arg struct {
	Key string
	Val any
}

type ContextLogger string
type ContextIdentifier string

const ctxLoggerKey ContextLogger = "ctx_logger"
const ctxIdentifier ContextIdentifier = "ctx_identifier"

func WrapLogger(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey, log)
}

func ExtractLogger(ctx context.Context) Logger {
	return ctx.Value(ctxLoggerKey).(Logger)
}

func WrapIdentifier(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxIdentifier, uuid.NewString())
}

func ExtractIdentifier(ctx context.Context) string {
	return ctx.Value(ctxIdentifier).(string)
}
