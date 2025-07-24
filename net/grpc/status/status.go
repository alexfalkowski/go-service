package status

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc/codes"
	"google.golang.org/grpc/status"
)

// Code is an alias for status.Code.
func Code(err error) codes.Code {
	return status.Code(err)
}

// Error is an alias for status.Error.
func Error(c codes.Code, msg string) error {
	return status.Error(c, msg)
}

// Errorf is an alias for status.Errorf.
func Errorf(c codes.Code, format string, a ...any) error {
	return status.Errorf(c, format, a...)
}
