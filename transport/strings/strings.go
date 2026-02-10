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

// IsIgnorable reports whether text contains one of the predefined ignorable substrings
// (for example "health", "metrics", etc.).
func IsIgnorable(text string) bool {
	return slices.ContainsFunc(ignorable, func(o string) bool { return strings.Contains(text, o) })
}

// IsFullMethod reports whether name is of the form `/package.service/method`.
func IsFullMethod(name string) bool {
	return strings.HasPrefix(name, "/") && strings.Count(name, "/") == 2 && strings.Count(name, ".") > 0
}

// SplitServiceMethod splits /package.service/method into package.service and method.
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
