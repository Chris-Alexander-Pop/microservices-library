package logger

import (
	"log/slog"
	"os"
)

type Config struct {
	Level  string `env:"LOG_LEVEL" env-default:"INFO"`
	Format string `env:"LOG_FORMAT" env-default:"JSON"` // JSON or TEXT
}

func New(cfg Config) *slog.Logger {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: parseLevel(cfg.Level),
	}

	if cfg.Format == "TEXT" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger) // Set as global default

	return logger
}

func parseLevel(level string) slog.Level {
	switch level {
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
