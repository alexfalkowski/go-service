package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/os"
)

func newTextLogger(params Params) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, handlerOptions(params.Config))).With(
		slog.String("id", params.ID.String()),
		slog.String("name", params.Name.String()),
		slog.String("version", params.Version.String()),
		slog.String("environment", params.Environment.String()),
	)
}
