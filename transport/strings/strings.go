package strings

import (
	"slices"

	"github.com/alexfalkowski/go-service/v2/strings"
)

var ignorable = []string{
	"health",
	"healthz",
	"livez",
	"readyz",
	"metrics",
	"favicon",
}

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

// Join is an alias for strings.Join.
func Join(sep string, ss ...string) string {
	return strings.Join(sep, ss...)
}

// IsEmpty is an alias for strings.IsEmpty.
func IsEmpty(s string) bool {
	return strings.IsEmpty(s)
}

// IsIgnorable reports whether text should be treated as ignorable by transport middleware.
//
// This helper is used across HTTP and gRPC transports to decide whether certain endpoints/methods should bypass
// middleware such as authentication, rate limiting, or logging.
//
// Matching behavior:
// IsIgnorable returns true when text contains any of the predefined ignorable substrings (for example "healthz"
// or "metrics"). This is intentionally a substring match (not an exact match), so callers should avoid passing
// overly broad inputs that could accidentally match unrelated paths/method names.
func IsIgnorable(text string) bool {
	return slices.ContainsFunc(ignorable, func(o string) bool { return strings.Contains(text, o) })
}

// IsFullMethod reports whether name is of the form `/package.service/method`.
//
// This is the canonical shape of gRPC full method names as they appear in interceptors (for example
// "/helloworld.Greeter/SayHello").
func IsFullMethod(name string) bool {
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
	return strings.Cut(name[1:], "/")
}

// ToLower is an alias for strings.ToLower.
func ToLower(s string) string {
	return strings.ToLower(s)
}
