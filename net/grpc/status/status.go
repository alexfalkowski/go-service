package status

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"google.golang.org/grpc/status"
)

// Status aliases the upstream gRPC status type.
//
// It is re-exported so callers can inspect gRPC errors through the go-service
// import path while preserving upstream behavior and methods such as Code,
// Message, Details, Err, and Proto.
type Status = status.Status

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

// FromError returns a status representation for err, if err is or wraps a gRPC
// status error.
//
// This is a thin wrapper around google.golang.org/grpc/status.FromError.
// Behavior is identical to the upstream implementation:
//
//   - If err was produced from a gRPC status, the returned Status reflects the
//     embedded code and message and ok is true.
//   - If err wraps a gRPC status error, the returned Status preserves the
//     underlying status while the message may incorporate wrapping context,
//     matching upstream behavior.
//   - Otherwise, ok is false and the returned Status represents codes.Unknown.
func FromError(err error) (*Status, bool) {
	return status.FromError(err)
}

// Error constructs a gRPC status error with code c and message msg.
//
// The returned error is suitable to be returned from a gRPC handler so the runtime can send the
// corresponding status code and message to the client. The message is client-visible by design.
//
// For structured status details (protobuf Any details), use the upstream status
// API directly (for example status.New(...).WithDetails(...)).
func Error(c codes.Code, msg string) error {
	return SafeError(c, msg, nil)
}

// SafeError wraps err with c and a message that is safe to send to clients.
//
// The wrapped error remains available through Unwrap for internal inspection, while gRPC sends msg instead
// of err.Error(). If c is codes.OK, SafeError returns nil to preserve upstream gRPC status invariants.
func SafeError(c codes.Code, msg string, err error) error {
	if c == codes.OK {
		return nil
	}

	return &statusError{code: c, msg: msg, err: err}
}

// Errorf constructs a formatted gRPC status error with code c.
//
// It formats the message using fmt-style formatting rules and returns an error suitable to be returned
// from a gRPC handler. The formatted message is client-visible by design.
//
// For structured status details (protobuf Any details), use the upstream status
// API directly (for example status.New(...).WithDetails(...)).
func Errorf(c codes.Code, format string, a ...any) error {
	return Error(c, fmt.Sprintf(format, a...))
}

type statusError struct {
	err  error
	msg  string
	code codes.Code
}

// Error returns the gRPC status error string with the safe message.
func (s *statusError) Error() string {
	return status.New(s.code, s.msg).Err().Error()
}

// SafeMessage returns the message that is safe to send to clients.
func (s *statusError) SafeMessage() string {
	return s.msg
}

// GRPCStatus returns the status sent by gRPC.
//
// The method name is part of google.golang.org/grpc/status.FromError's contract: gRPC detects custom
// status errors by looking for GRPCStatus() *Status.
func (s *statusError) GRPCStatus() *Status {
	return status.New(s.code, s.msg)
}

// Is reports whether target has the same gRPC status code and message.
func (s *statusError) Is(target error) bool {
	type grpcStatus interface {
		GRPCStatus() *Status
	}

	if target, ok := target.(grpcStatus); ok {
		t := target.GRPCStatus()
		return t != nil && s.code == t.Code() && s.msg == t.Message()
	}

	return false
}

// Unwrap returns the wrapped cause, if any.
func (s *statusError) Unwrap() error {
	return s.err
}
