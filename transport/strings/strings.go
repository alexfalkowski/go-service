package strings

import (
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
	for _, o := range observables {
		if strings.Contains(text, o) {
			return true
		}
	}

	return false
}
