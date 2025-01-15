package header

import (
	"github.com/alexfalkowski/go-service/os"
)

// Map for tracer.
type Map map[string]string

// Secrets will traverse the map and load any secrets that have been configured.
func (m Map) Secrets() error {
	for k, v := range m {
		if !os.FileExists(v) {
			continue
		}

		m[k] = os.MustReadFile(v)
	}

	return nil
}
