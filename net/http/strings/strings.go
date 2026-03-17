package strings

import (
	"slices"

	"github.com/alexfalkowski/go-service/v2/strings"
)

var ignorablePaths = []string{
	"healthz",
	"livez",
	"readyz",
	"metrics",
	"favicon",
	"favicon.ico",
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

// IsEmpty is an alias for strings.IsEmpty.
func IsEmpty(s string) bool {
	return strings.IsEmpty(s)
}

// Join is an alias for strings.Join.
func Join(sep string, ss ...string) string {
	return strings.Join(sep, ss...)
}

// IsIgnorable reports whether path should be treated as ignorable by HTTP middleware.
//
// Matching is exact by path segment so only well-known operational endpoints such as
// `/<service>/healthz` or `/<service>/metrics` are ignored.
func IsIgnorable(path string) bool {
	path = strings.Trim(path, "/")
	if strings.IsEmpty(path) {
		return false
	}
	if strings.Count(path, "/") > 1 {
		return false
	}

	last := path
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		last = path[idx+1:]
	}

	return slices.Contains(ignorablePaths, last)
}

// ToLower is an alias for strings.ToLower.
func ToLower(s string) string {
	return strings.ToLower(s)
}
