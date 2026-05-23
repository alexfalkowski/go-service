package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/strings"
)

func handlerOptions(cfg *Config) *slog.HandlerOptions {
	return &slog.HandlerOptions{
		Level: level(cfg),
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.LevelKey {
				if level, ok := attr.Value.Any().(slog.Level); ok {
					attr.Value = slog.StringValue(strings.ToLower(level.String()))
				}
			}

			return attr
		},
	}
}
