package logger

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/telemetry/logger"
)

// LevelError is an alias of [github.com/alexfalkowski/go-service/v2/telemetry/logger.LevelError].
const LevelError = logger.LevelError

// Logger is an alias of [github.com/alexfalkowski/go-service/v2/telemetry/logger.Logger].
//
// It is re-exported here so transport-layer code can depend on a single logger type when composing
// interceptors.
type Logger = logger.Logger

// CodeToLevel translates a gRPC status code to a [github.com/alexfalkowski/go-service/v2/telemetry/logger.Level].
//
// The mapping is intentionally coarse-grained:
//
//   - [codes.OK] is logged at info.
//   - Client/expected error codes (cancellation, invalid argument, not found, unauthenticated, etc.) are logged at warn.
//   - All other codes (typically server/transient failures) are logged at error.
func CodeToLevel(code codes.Code) logger.Level {
	switch code {
	case codes.OK:
		return logger.LevelInfo
	case codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated,
		codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange:
		return logger.LevelWarn
	default:
		return logger.LevelError
	}
}

func message(msg string) string {
	return "grpc: " + msg
}
