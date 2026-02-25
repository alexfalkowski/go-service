package logger

import "github.com/alexfalkowski/go-service/v2/di"

// Module wires logger construction into Fx.
//
// It provides:
//   - `NewLogger`, which constructs the configured logger and installs it as the
//     process-wide slog default.
//   - `convertLogger`, which exposes the underlying `*slog.Logger` for packages
//     that depend directly on slog.
//
// When logging is disabled (`*Config` is nil), `NewLogger` returns a nil `*Logger`
// and `convertLogger` will therefore provide a nil `*slog.Logger`.
var Module = di.Module(
	di.Constructor(NewLogger),
	di.Constructor(convertLogger),
)
