package strings

import (
	"slices"
	"strings"
)

var observables = []string{
	"health",
	"healthz",
	"livez",
	"readyz",
	"metrics",
}

// IsObservable in the text.
func IsObservable(text string) bool {
	return slices.ContainsFunc(observables, func(o string) bool { return strings.Contains(text, o) })
}
