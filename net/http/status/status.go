package status

import (
	"fmt"

	"github.com/alexfalkowski/go-service/v2/context"
	"github.com/alexfalkowski/go-service/v2/errors"
	"github.com/alexfalkowski/go-service/v2/net/grpc/status"
	"github.com/alexfalkowski/go-service/v2/net/http"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// Error constructs an error that carries an HTTP status code and a message.
//
// The returned error implements the Coder interface so handlers and helpers can extract the HTTP
// status code via Code(err). The message is considered safe to send to clients by WriteError.
func Error(code int, msg string) error {
	return &statusError{code: code, error: msg, msg: msg}
}

// SafeError wraps err with code and a safe HTTP-prefixed status message that is safe to send to clients.
//
// The wrapped error remains available through Unwrap for internal inspection, while WriteError sends the safe message
// instead of err.Error(). If err already carries a status code, it is returned unchanged.
func SafeError(code int, err error) error {
	code = normalizeCode(code, err)
	msg := DefaultMessage(code)

	if err == nil {
		return &statusError{code: code, error: msg, msg: msg}
	}

	if IsError(err) {
		return err
	}

	return &statusError{code: code, error: err.Error(), msg: msg, err: err}
}

// InternalServerError wraps err with StatusInternalServerError (500) unless err already carries a status code.
//
// This is a convenience wrapper over FromError.
func InternalServerError(err error) error {
	return FromError(http.StatusInternalServerError, err)
}

// ServiceUnavailableError wraps err with StatusServiceUnavailable (503) unless err already carries a status code.
//
// This is a convenience wrapper over FromError.
func ServiceUnavailableError(err error) error {
	return FromError(http.StatusServiceUnavailable, err)
}

// UnauthorizedError wraps err with StatusUnauthorized (401) unless err already carries a status code.
//
// This is a convenience wrapper over FromError.
func UnauthorizedError(err error) error {
	return FromError(http.StatusUnauthorized, err)
}

// BadRequestError wraps err with StatusBadRequest (400) unless err already carries a status code.
//
// This is a convenience wrapper over FromError.
func BadRequestError(err error) error {
	return FromError(http.StatusBadRequest, err)
}

// FromError returns err unchanged if it already carries a status code; otherwise it wraps err with code.
//
// This helper is intentionally idempotent for errors already produced by this package (or any error
// implementing Coder): calling FromError on such errors does not overwrite the original status code.
// Raw *[http.MaxBytesError] values are normalized to StatusRequestEntityTooLarge (413) regardless of the
// provided code so oversized request bodies surface consistently.
// The wrapped error message remains diagnostic through Error, but WriteError sends a safe status message instead
// of err.Error().
//
// Passing nil returns a status error with the default safe message for code.
func FromError(code int, err error) error {
	return SafeError(code, err)
}

// Errorf formats a message and returns an error with the provided status code.
//
// This is a convenience wrapper over Error(code, [fmt.Sprintf](...)).
func Errorf(code int, format string, a ...any) error {
	return Error(code, fmt.Sprintf(format, a...))
}

// IsError reports whether err carries a status code.
//
// It returns true for:
//   - errors produced by this package, and
//   - any error implementing the Coder interface.
func IsError(err error) bool {
	_, ok := coderFromError(err)
	return ok
}

// DefaultMessage returns the default safe message for code.
//
// Unknown codes fall back to StatusInternalServerError.
func DefaultMessage(code int) string {
	if code == http.StatusClientClosedRequest {
		return "http: client closed request"
	}

	if msg := http.StatusText(code); msg != "" {
		return "http: " + strings.ToLower(msg)
	}

	return "http: " + strings.ToLower(http.StatusText(http.StatusInternalServerError))
}

// Code extracts the HTTP status code from err.
//
// Resolution order:
//  1. If err implements Coder, return coder.Code().
//  2. If err is a raw *[http.MaxBytesError], return StatusRequestEntityTooLarge (413).
//  3. If err wraps context.Canceled, return StatusClientClosedRequest (499).
//  4. If err wraps context.DeadlineExceeded, return StatusGatewayTimeout (504).
//  5. If err is a gRPC status error, map its gRPC code to an HTTP status code using the statusCodes table.
//  6. Otherwise return StatusInternalServerError (500).
func Code(err error) int {
	if coder, ok := coderFromError(err); ok {
		return coder.Code()
	}

	if isMaxBytesError(err) {
		return http.StatusRequestEntityTooLarge
	}

	if errors.Is(err, context.Canceled) {
		return http.StatusClientClosedRequest
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return http.StatusGatewayTimeout
	}

	s, ok := status.FromError(err)
	if ok {
		if code, ok := statusCodes[s.Code()]; ok {
			return code
		}

		return http.StatusInternalServerError
	}

	return http.StatusInternalServerError
}

func coderFromError(err error) (Coder, bool) {
	var coder Coder
	if errors.As(err, &coder) {
		return coder, true
	}

	return nil, false
}

func isMaxBytesError(err error) bool {
	_, ok := errors.AsType[*http.MaxBytesError](err)
	return ok
}

func normalizeCode(code int, err error) int {
	if isMaxBytesError(err) {
		return http.StatusRequestEntityTooLarge
	}

	return code
}

type statusError struct {
	err   error
	error string
	msg   string
	code  int
}

// Code returns the HTTP status code carried by the error.
func (s *statusError) Code() int {
	return s.code
}

// Error returns the diagnostic error message.
func (s *statusError) Error() string {
	return s.error
}

// SafeMessage returns the status message that is safe to send to clients.
func (s *statusError) SafeMessage() string {
	return s.msg
}

// Unwrap returns the wrapped cause, if any.
func (s *statusError) Unwrap() error {
	return s.err
}
