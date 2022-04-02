package strings

import (
	"regexp"
)

var words = regexp.MustCompile(`health|liveness|readiness`)

// IsHealth in the text.
func IsHealth(text string) bool {
	return words.MatchString(text)
}
