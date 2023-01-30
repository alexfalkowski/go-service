package strings

import (
	"regexp"
)

var health = regexp.MustCompile(`health|healthz|livez|readyz|metrics`)

// IsHealth in the text.
func IsHealth(text string) bool {
	return health.MatchString(text)
}
