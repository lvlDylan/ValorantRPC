package logging

import (
	"ValorantRPC/pkg/errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type Config struct {
	Level      string
	JSONFormat bool
	FilePath   string
}

type Log struct {
	logger *slog.Logger
}

func parseLevel(levelStr string) slog.Level {
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func New(cfg Config) (*Log, error) {
	var writer io.Writer = os.Stdout

	if cfg.FilePath != "" {
		file, err := os.OpenFile(cfg.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			cause := fmt.Errorf("%w : %v", errors.ErrLoggerCannotOpenFile, err)
			return nil, errors.New("logging.New", cause)
		}
		writer = file
	}

	opts := &slog.HandlerOptions{
		Level: parseLevel(cfg.Level),
	}

	var handler slog.Handler
	if cfg.JSONFormat {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	slogLogger := slog.New(handler)
	return &Log{
		logger: slogLogger,
	}, nil
}

func (l *Log) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *Log) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *Log) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *Log) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}
