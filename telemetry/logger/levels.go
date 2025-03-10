package logger

import "log/slog"

var levels = map[string]slog.Level{
	"debug": slog.LevelDebug,
	"info":  slog.LevelInfo,
	"warn":  slog.LevelWarn,
	"error": slog.LevelError,
}

func level(cfg *Config) slog.Level {
	return levels[cfg.Level]
}
