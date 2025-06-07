package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/lmittmann/tint"
)

func newTintLogger(params LoggerParams) *slog.Logger {
	opts := &tint.Options{
		Level: level(params.Config),
		ReplaceAttr: func(groups []string, attr slog.Attr) slog.Attr {
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

	return slog.New(tint.NewHandler(os.Stdout, opts)).With(
		slog.String("id", params.ID.String()),
		slog.String("name", params.Name.String()),
		slog.String("version", params.Version.String()),
		slog.String("environment", params.Environment.String()),
	)
}
