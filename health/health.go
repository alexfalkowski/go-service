package health

import (
	"regexp"
)

var words = regexp.MustCompile(`health|liveness|readiness`)

// Is part of health.
func Is(text string) bool {
	return words.MatchString(text)
}
