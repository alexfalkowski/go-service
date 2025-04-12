package strings

import "strings"

var (
	// Contains is an alias for strings.Contains.
	Contains = strings.Contains

	// Cut is an alias for strings.Cut.
	Cut = strings.Cut

	// HasPrefix is an alias for strings.HasPrefix.
	HasPrefix = strings.HasPrefix

	// Repeat is an alias for strings.Repeat.
	Repeat = strings.Repeat

	// ToLower is an alias for strings.ToLower.
	ToLower = strings.ToLower

	// TrimSpace is an alias for strings.TrimSpace.
	TrimSpace = strings.TrimSpace
)

// IsEmpty checks if a string is empty.
func IsEmpty(s string) bool {
	return len(s) == 0
}
