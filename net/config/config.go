package config

import "github.com/alexfalkowski/go-service/v2/time"

// Duration returns the duration from the options or a timeout value.
func Duration(options map[string]string, key string, timeout time.Duration) time.Duration {
	if val, ok := options[key]; ok {
		return time.MustParseDuration(val)
	}
	return timeout
}
