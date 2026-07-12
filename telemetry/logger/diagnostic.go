package logger

import (
	"log/slog"

	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/lmittmann/tint"
)

// NewDiagnosticLogger builds a stdout logger for internal diagnostics that
// mirrors the configured logger format without depending on the OTLP export
// pipeline.
//
// It selects the stdout handler for cfg.Kind ("json", "text", or "tint") and
// applies the configured level. The "otlp" kind, an unknown kind, and a nil cfg
// (logging disabled) all fall back to JSON on stdout, because those have no
// local stdout format and this sink must never route through OTLP.
//
// Unlike the application logger, the returned logger does not attach identity
// resource attributes or trace context; it is intended for a self-contained
// local diagnostic sink such as the OpenTelemetry error handler.
func NewDiagnosticLogger(cfg *Config) *slog.Logger {
	switch cfg.GetKind() {
	case "text":
		return slog.New(slog.NewTextHandler(os.Stdout, handlerOptions(cfg)))
	case "tint":
		return slog.New(tint.NewTextHandler(os.Stdout, tintOptions(cfg)))
	default:
		return slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions(cfg)))
	}
}
