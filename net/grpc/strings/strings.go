package strings

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc/health"
	"github.com/alexfalkowski/go-service/v2/net/url"
	"github.com/alexfalkowski/go-service/v2/strings"
)

const (
	// Empty is an alias for strings.Empty.
	Empty = strings.Empty

	// Space is an alias for strings.Space.
	Space = strings.Space
)

// Bytes is an alias for strings.Bytes.
func Bytes(s string) []byte {
	return strings.Bytes(s)
}

// Contains is an alias for strings.Contains.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// HasPrefix is an alias for strings.HasPrefix.
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// Cut is an alias for strings.Cut.
func Cut(s, sep string) (string, string, bool) {
	return strings.Cut(s, sep)
}

// IsEmpty is an alias for strings.IsEmpty.
func IsEmpty(s string) bool {
	return strings.IsEmpty(s)
}

// Join is an alias for strings.Join.
func Join(sep string, ss ...string) string {
	return strings.Join(sep, ss...)
}

// IsOperationMethod reports whether name is a gRPC operation method owned by the transport.
//
// Matching is exact for the standard gRPC health service methods.
func IsOperationMethod(name string) bool {
	switch name {
	case health.CheckFullMethodName, health.ListFullMethodName, health.WatchFullMethodName:
		return true
	default:
		return false
	}
}

// IsFullMethod reports whether name is of the form `/package.service/method`.
//
// This is the canonical shape of gRPC full method names as they appear in interceptors (for example
// "/helloworld.Greeter/SayHello").
func IsFullMethod(name string) bool {
	// Buf-managed protos in this repo require a package, so service names are
	// expected to include a package-qualified dot, e.g. "/greet.v1.Greeter/SayHello".
	return strings.HasPrefix(name, "/") && strings.Count(name, "/") == 2 && strings.Count(name, ".") > 0
}

// SplitServiceMethod splits a gRPC full method name into service and method components.
//
// If name is not a valid gRPC full method (see IsFullMethod), it returns ("", "", false).
// Otherwise it returns ("package.service", "method", true).
func SplitServiceMethod(name string) (string, string, bool) {
	if !IsFullMethod(name) {
		return "", "", false
	}
	return url.SplitPath(name)
}
