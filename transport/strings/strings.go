package strings

import (
	"regexp"
)

var health = regexp.MustCompile(`health|healthz|livez|readyz|metrics`)

// IsObservable in the text.
func IsObservable(text string) bool {
	return health.MatchString(text)
}
