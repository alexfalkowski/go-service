package logger

import (
	"log/slog"
	"os"
	"strings"
)

func newStdoutLogger(params Params) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: level(params.Config),
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Key == slog.LevelKey {
				level := attr.Value.Any().(slog.Level)
				attr.Value = slog.StringValue(strings.ToLower(level.String()))
			}

			return attr
		},
	}

	var handler slog.Handler

	if params.Environment.IsDevelopment() {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	return slog.New(handler).With(
		slog.String("name", params.Name.String()),
		slog.String("version", params.Version.String()),
		slog.String("environment", params.Environment.String()),
	)
}
