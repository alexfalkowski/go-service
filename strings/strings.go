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
func CutColon(arg string) (string, string) {
	before, after, ok := strings.Cut(arg, ":")
	if !ok {
		return "", ""
	}

	return before, after
}
