package logger

import "github.com/alexfalkowski/go-signal/v2"

// Module wires logger construction into [github.com/alexfalkowski/go-signal/v2].
//
// It provides:
//   - [NewLogger], which constructs the configured logger and installs it as the
//     process-wide slog default.
//   - `convertLogger`, which exposes the underlying `*slog.Logger` for packages
//     that depend directly on slog.
//
// When logging is disabled (*[Config] is nil), [NewLogger] returns a nil *[Logger]
// and `convertLogger` will therefore provide a nil `*slog.Logger`.
var Module = signal.Module(
	signal.Constructor(NewLogger),
	signal.Constructor(convertLogger),
)
