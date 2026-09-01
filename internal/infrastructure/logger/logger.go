package logger

import (
	"context"
	"time"
)

type Logger interface {
	Debug(msg string, args ...Field)
	Info(msg string, args ...Field)
	Warn(msg string, args ...Field)
	Error(msg string, args ...Field)
	Fatal(msg string, args ...Field)
	With(args ...Field) Logger
	Close()
}

type loggerContextKey struct{}

var key = loggerContextKey{}

func WithContext(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, key, log)
}

func FromContext(ctx context.Context) Logger {
	log, ok := ctx.Value(key).(Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

type Field struct {
	Key   string
	Value any
}

func String(key string, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Duration(key string, value interface{ String() string }) Field {
	return Field{Key: key, Value: value}
}

func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value}
}

func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err.Error()}
}
