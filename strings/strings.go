package strings

import (
	"regexp"
)

var health = regexp.MustCompile(`health|liveness|readiness|metrics`)

// IsHealth in the text.
func IsHealth(text string) bool {
	return health.MatchString(text)
}
