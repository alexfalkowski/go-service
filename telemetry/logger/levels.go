package logger

import "log/slog"

var levels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func validateLevel(cfg *Config) error {
	if cfg.Level == "" {
		return nil
	}
	if _, ok := levels[cfg.Level]; ok {
		return nil
	}

	return ErrInvalidLevel
}

func level(cfg *Config) slog.Level {
	if level, ok := levels[cfg.Level]; ok {
		return level
	}

	return slog.LevelInfo
}
