package options

import "github.com/alexfalkowski/go-service/v2/time"

// Map contains key-value pairs.
type Map map[string]string

// Duration returns the duration from the options or a timeout value.
func (m Map) Duration(key string, timeout time.Duration) time.Duration {
	if val, ok := m[key]; ok {
		return time.MustParseDuration(val)
	}
	return timeout
}
