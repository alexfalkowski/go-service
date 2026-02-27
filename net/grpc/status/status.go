package status

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"google.golang.org/grpc/status"
)

// Code returns the gRPC status code for err as a go-service codes.Code.
//
// This is a thin wrapper around google.golang.org/grpc/status.Code, returning
// the code as net/grpc/codes.Code (which aliases the upstream codes.Code type).
//
// Behavior is defined by the upstream gRPC implementation. In particular:
//
//   - If err is nil, the returned code is codes.OK.
//   - If err is a gRPC status error, the returned code is the status code
//     contained in that error.
//   - If err is not a gRPC status error, upstream behavior typically returns
//     codes.Unknown.
//
// Use this function in clients/handlers to classify failures and decide on a
// response strategy (retry, map to HTTP status codes, etc.).
func Code(err error) codes.Code {
	return status.Code(err)
}

// Error constructs a gRPC status error with code c and message msg.
//
// This is a thin wrapper around google.golang.org/grpc/status.Error. The
// returned error is suitable to be returned from a gRPC handler so the runtime
// can send the corresponding status code and message to the client.
//
// For structured status details (protobuf Any details), use the upstream status
// API directly (for example status.New(...).WithDetails(...)).
func Error(c codes.Code, msg string) error {
	return status.Error(c, msg)
}

// Errorf constructs a formatted gRPC status error with code c.
//
// This is a thin wrapper around google.golang.org/grpc/status.Errorf. It formats
// the message using fmt-style formatting rules and returns an error suitable to
// be returned from a gRPC handler.
//
// For structured status details (protobuf Any details), use the upstream status
// API directly (for example status.New(...).WithDetails(...)).
func Errorf(c codes.Code, format string, a ...any) error {
	return status.Errorf(c, format, a...)
}
