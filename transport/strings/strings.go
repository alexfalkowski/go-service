package strings

import (
	"slices"

	"github.com/alexfalkowski/go-service/v2/strings"
)

var observables = []string{
	"health",
	"healthz",
	"livez",
	"readyz",
	"metrics",
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

// IsObservable in the text.
func IsObservable(text string) bool {
	return slices.ContainsFunc(observables, func(o string) bool { return strings.Contains(text, o) })
}

// SplitServiceMethod will split /package.service/method to package.service and method.
func SplitServiceMethod(name string) (string, string) {
	service, method, _ := strings.Cut(name[1:], "/")
	return service, method
}

// ToLower is an alias for strings.ToLower.
func ToLower(s string) string {
	return strings.ToLower(s)
}
