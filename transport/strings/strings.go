package strings

import (
	"slices"

	"github.com/alexfalkowski/go-service/v2/strings"
)

var (
	observables = []string{
		"health",
		"healthz",
		"livez",
		"readyz",
		"metrics",
	}

	// Bytes is an alias for strings.Bytes.
	Bytes = strings.Bytes

	// Cut is an alias for strings.Cut.
	Cut = strings.Cut

	// Join is an alias for strings.Join.
	Join = strings.Join

	// IsEmpty is an alias for strings.IsEmpty.
	IsEmpty = strings.IsEmpty

	// ToLower is an alias for strings.ToLower.
	ToLower = strings.ToLower
)

// IsObservable in the text.
func IsObservable(text string) bool {
	return slices.ContainsFunc(observables, func(o string) bool { return strings.Contains(text, o) })
}
