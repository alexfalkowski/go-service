package strings

import (
	"slices"
	"strings"
)

var (
	// Contains is an alias for strings.Contains.
	Contains = strings.Contains

	// Cut is an alias for strings.Cut.
	Cut = strings.Cut

	// HasPrefix is an alias for strings.HasPrefix.
	HasPrefix = strings.HasPrefix

	// Repeat is an alias for strings.Repeat.
	Repeat = strings.Repeat

	// ReplaceAll is an alias for strings.ReplaceAll.
	ReplaceAll = strings.ReplaceAll

	// ToLower is an alias for strings.ToLower.
	ToLower = strings.ToLower

	// TrimSpace is an alias for strings.TrimSpace.
	TrimSpace = strings.TrimSpace
)

// IsEmpty checks if a string is empty.
func IsEmpty(s string) bool {
	return len(s) == 0
}

// IsAnyEmpty checks if any strings are empty.
func IsAnyEmpty(ss ...string) bool {
	return slices.ContainsFunc(ss, IsEmpty)
}

// Join strings by a separator.
// This allows to do strings.Join(" ", "1", "2").
func Join(sep string, ss ...string) string {
	return strings.Join(ss, sep)
}

// Concat will take all the strings and join them with an empty string.
func Concat(ss ...string) string {
	return Join("", ss...)
}

// CutColon will split by : and give before and after.
func CutColon(s string) (string, string) {
	before, after, _ := Cut(s, ":")

	return before, after
}
