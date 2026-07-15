package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/strings"
)

var levels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func validateLevel(cfg *Config) error {
	if strings.IsEmpty(cfg.Level) {
		return nil
	}
	if _, ok := levels[cfg.Level]; ok {
		return nil
	}

	return ErrInvalidLevel
}

func level(cfg *Config) slog.Level {
	if cfg == nil {
		return slog.LevelInfo
	}
	if level, ok := levels[cfg.Level]; ok {
		return level
	}

	return slog.LevelInfo
}
