package logger

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
