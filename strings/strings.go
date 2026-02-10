package strings

import (
	"slices"
	"strings"
)

const (
	// Empty is the empty string constant.
	Empty = ""

	// Space is the space string constant.
	Space = " "
)

// Contains is an alias for strings.Contains.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// Count is an alias for strings.Count.
func Count(s, substr string) int {
	return strings.Count(s, substr)
}

// Cut is an alias for strings.Cut.
func Cut(s, sep string) (string, string, bool) {
	return strings.Cut(s, sep)
}

// HasPrefix is an alias for strings.HasPrefix.
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// Repeat is an alias for strings.Repeat.
func Repeat(s string, count int) string {
	return strings.Repeat(s, count)
}

// ReplaceAll is an alias for strings.ReplaceAll.
func ReplaceAll(s, o, n string) string {
	return strings.ReplaceAll(s, o, n)
}

// ToLower is an alias for strings.ToLower.
func ToLower(s string) string {
	return strings.ToLower(s)
}

// TrimSpace is an alias for strings.TrimSpace.
func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}

// IsEmpty reports whether s is empty.
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsAnyEmpty reports whether any of ss are empty.
func IsAnyEmpty(ss ...string) bool {
	return slices.ContainsFunc(ss, IsEmpty)
}

// Join joins ss with sep.
//
// This helper allows joining variadic strings without allocating a slice at the callsite.
func Join(sep string, ss ...string) string {
	return strings.Join(ss, sep)
}

// Concat concatenates ss without a separator.
func Concat(ss ...string) string {
	return Join(Empty, ss...)
}

// CutColon splits s on the first ":" and returns the parts before and after.
//
// If ":" is not present, the returned after value is empty.
func CutColon(s string) (string, string) {
	before, after, _ := Cut(s, ":")
	return before, after
}
