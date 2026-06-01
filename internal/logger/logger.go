package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"

	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level      string // "DEBUG", "INFO", "WARN", "ERROR"
	Format     string // "text" (con tint, para dev) o "json" (para prod)
	File       string // path al archivo; vacío = solo stdout
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
}

func New(cfg Config) *slog.Logger {
	level := parseLevel(cfg.Level)

	writer := buildWriter(cfg)

	var baseHandler slog.Handler
	if cfg.Format == "text" {
		baseHandler = tint.NewHandler(writer, &tint.Options{
			Level:      level,
			TimeFormat: time.TimeOnly,
			AddSource:  false,
		})
	} else {
		baseHandler = slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level:     level,
			AddSource: false,
		})
	}
	handler := &ContextHandler{Handler: baseHandler}

	return slog.New(handler)
}

func buildWriter(cfg Config) io.Writer {
	if cfg.File == "" {
		return os.Stdout
	}

	fileWriter := &lumberjack.Logger{
		Filename:   cfg.File,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	}

	return io.MultiWriter(os.Stdout, fileWriter)
}

func parseLevel(s string) slog.Level {

	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN", "WARNING":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func SetDefault(l *slog.Logger) {
	slog.SetDefault(l)
}
