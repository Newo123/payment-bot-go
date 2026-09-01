package zap

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Newo123/payment-bot-go/internal/infrastructure/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	file *os.File
}

func NewLogger(config Config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		config.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(2))

	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

func (l *Logger) Debug(msg string, args ...logger.Field) {
	l.log(msg, zap.DebugLevel, args...)
}
func (l *Logger) Info(msg string, args ...logger.Field) {
	l.log(msg, zap.InfoLevel, args...)
}
func (l *Logger) Warn(msg string, args ...logger.Field) {
	l.log(msg, zap.WarnLevel, args...)
}
func (l *Logger) Error(msg string, args ...logger.Field) {
	l.log(msg, zap.ErrorLevel, args...)
}
func (l *Logger) Fatal(msg string, args ...logger.Field) {
	l.log(msg, zap.FatalLevel, args...)
	os.Exit(1)
}

func (l *Logger) With(fields ...logger.Field) logger.Logger {
	zapFields := make([]zap.Field, 0, len(fields))

	for _, field := range fields {
		zapFields = append(zapFields, zap.Any(field.Key, field.Value))
	}

	return &Logger{
		Logger: l.Logger.With(zapFields...),
		file:   l.file,
	}
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err)
	}
}

func (l *Logger) log(msg string, level zapcore.Level, fields ...logger.Field) {
	if ce := l.Logger.Check(level, msg); ce != nil {
		zfields := make([]zap.Field, 0, len(fields))

		for _, f := range fields {
			zfields = append(zfields, zap.Any(f.Key, f.Value))
		}

		ce.Write(zfields...)
	}
}
