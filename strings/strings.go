package strings

import (
	"slices"
	"strings"
)

const (
	// Empty for strings.
	Empty = ""

	// Space for strings.
	Space = " "
)

// Contains is an alias for strings.Contains.
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
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

// IsEmpty checks if a string is empty.
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsAnyEmpty checks if any strings are empty.
func IsAnyEmpty(ss ...string) bool {
	return slices.ContainsFunc(ss, IsEmpty)
}

// Join strings by a separator.
// This allows to do strings.Join(strings.Space, "1", "2").
func Join(sep string, ss ...string) string {
	return strings.Join(ss, sep)
}

// Concat will take all the strings and join them with an empty string.
func Concat(ss ...string) string {
	return Join(Empty, ss...)
}

// CutColon will split by : and give before and after.
func CutColon(s string) (string, string) {
	before, after, _ := Cut(s, ":")
	return before, after
}
