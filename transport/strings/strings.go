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

	// Join is an alias for strings.Join.
	Join = strings.Join
)

// IsObservable in the text.
func IsObservable(text string) bool {
	return slices.ContainsFunc(observables, func(o string) bool { return strings.Contains(text, o) })
}
