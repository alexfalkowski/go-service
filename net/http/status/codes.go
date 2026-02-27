package status

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"github.com/alexfalkowski/go-service/v2/net/http"
)

// statusCodes maps gRPC status codes to HTTP status codes.
//
// This table is used by Code(err) when err is a gRPC status error (i.e. status.FromError(err) succeeds).
// It provides a reasonable default mapping for exposing gRPC-style failures over HTTP.
//
// Provenance:
// The mapping is based on the grpc-gateway runtime mapping:
// https://github.com/grpc-ecosystem/grpc-gateway/blob/main/runtime/errors.go
//
// Notes:
//   - gRPC Canceled is mapped to 499 (Client Closed Request), which is a non-standard but commonly used
//     code to represent client cancellations in HTTP environments.
var statusCodes = map[codes.Code]int{
	codes.OK:                 http.StatusOK,
	codes.Canceled:           499,
	codes.Unknown:            http.StatusInternalServerError,
	codes.InvalidArgument:    http.StatusBadRequest,
	codes.DeadlineExceeded:   http.StatusGatewayTimeout,
	codes.NotFound:           http.StatusNotFound,
	codes.AlreadyExists:      http.StatusConflict,
	codes.PermissionDenied:   http.StatusForbidden,
	codes.Unauthenticated:    http.StatusUnauthorized,
	codes.ResourceExhausted:  http.StatusTooManyRequests,
	codes.FailedPrecondition: http.StatusBadRequest,
	codes.Aborted:            http.StatusConflict,
	codes.OutOfRange:         http.StatusBadRequest,
	codes.Unimplemented:      http.StatusNotImplemented,
	codes.Internal:           http.StatusInternalServerError,
	codes.Unavailable:        http.StatusInternalServerError,
	codes.DataLoss:           http.StatusInternalServerError,
}
