package strings

import (
	"slices"
	"strings"
)

const (
	// Empty is the empty string constant.
	//
	// It is provided as a named constant for readability and reuse.
	Empty = ""

	// Space is the single ASCII space character (" ").
	//
	// It is provided as a named constant for readability and reuse.
	Space = " "
)

// Contains reports whether substr is within s.
//
// This is a thin wrapper around strings.Contains and does not change semantics.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Count counts the number of non-overlapping instances of substr in s.
//
// This is a thin wrapper around strings.Count and does not change semantics.
func Count(s, substr string) int {
	return strings.Count(s, substr)
}

// Cut slices s around the first instance of sep, returning the text before and after sep.
//
// The returned boolean reports whether sep was found.
//
// This is a thin wrapper around strings.Cut and does not change semantics.
func Cut(s, sep string) (string, string, bool) {
	return strings.Cut(s, sep)
}

// HasPrefix reports whether s begins with prefix.
//
// This is a thin wrapper around strings.HasPrefix and does not change semantics.
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// Repeat returns a new string consisting of count copies of s.
//
// This is a thin wrapper around strings.Repeat and does not change semantics.
func Repeat(s string, count int) string {
	return strings.Repeat(s, count)
}

// ReplaceAll returns a copy of s with all non-overlapping instances of o replaced by n.
//
// This is a thin wrapper around strings.ReplaceAll and does not change semantics.
func ReplaceAll(s, o, n string) string {
	return strings.ReplaceAll(s, o, n)
}

// ToLower returns s with all Unicode letters mapped to their lower case.
//
// This is a thin wrapper around strings.ToLower and does not change semantics.
func ToLower(s string) string {
	return strings.ToLower(s)
}

// TrimSpace returns s with all leading and trailing white space removed.
//
// White space is defined by Unicode.
//
// This is a thin wrapper around strings.TrimSpace and does not change semantics.
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty reports whether s is empty.
//
// This is equivalent to len(s) == 0, but provided for readability at call sites.
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsAnyEmpty reports whether any of ss are empty.
//
// It returns true when at least one element of ss has length 0.
func IsAnyEmpty(ss ...string) bool {
	return slices.ContainsFunc(ss, IsEmpty)
}

// Join joins ss with sep.
//
// This is equivalent to strings.Join, but it accepts variadic input so callers can
// avoid allocating a slice at the callsite when they already have discrete string
// values.
func Join(sep string, ss ...string) string {
	return strings.Join(ss, sep)
}

// Concat concatenates ss without a separator.
//
// This is equivalent to Join(Empty, ss...).
func Concat(ss ...string) string {
	return Join(Empty, ss...)
}

// CutColon splits s on the first ":" and returns the parts before and after.
//
// If ":" is not present, the returned after value is empty.
//
// This helper is commonly used by go-service "source string" conventions where a
// value is prefixed with a kind such as "env:NAME" or "file:/path".
func CutColon(s string) (string, string) {
	before, after, _ := Cut(s, ":")
	return before, after
}
