package strings

import (
	"github.com/alexfalkowski/go-service/v2/net/grpc/health"
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
