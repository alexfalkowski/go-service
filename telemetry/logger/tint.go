package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/lmittmann/tint"
)

func tintOptions(cfg *Config) *tint.Options {
	return &tint.Options{
		Level: level(cfg),
		ReplaceAttr: func(_ []string, attr slog.Attr) slog.Attr {
			if attr.Value.Kind() == slog.KindAny {
				if err, ok := attr.Value.Any().(error); ok {
					err := tint.Err(err)
					err.Key = attr.Key

					return err
				}
			}

			return attr
		},
	}
}

func newTintLogger(params LoggerParams) *slog.Logger {
	return slog.New(NewTraceHandler(tint.NewTextHandler(os.Stdout, tintOptions(params.Config)))).With(
		slog.String("id", params.ID.String()),
		slog.String("name", params.Name.String()),
		slog.String("version", params.Version.String()),
		slog.String("environment", params.Environment.String()),
	)
}
