package strings

import (
	"github.com/alexfalkowski/go-service/v2/env"
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

// Concat is an alias for strings.Concat.
func Concat(ss ...string) string {
	return strings.Concat(ss...)
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

// IsOperationPath reports whether path is a service-owned operational endpoint.
//
// Matching is exact so application routes such as "/admin/metrics" are not
// treated as transport operation endpoints for auth or rate-limit bypasses.
func IsOperationPath(name env.Name, path string) bool {
	if strings.IsEmpty(path) || path[0] != '/' {
		return false
	}

	service, operation, ok := strings.Cut(path[1:], "/")
	if !ok || service != name.String() {
		return false
	}

	switch operation {
	case "healthz", "livez", "readyz", "metrics":
		return true
	default:
		return false
	}
}

// ToLower is an alias for strings.ToLower.
func ToLower(s string) string {
	return strings.ToLower(s)
}
