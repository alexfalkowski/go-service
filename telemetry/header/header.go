package header

import (
	"github.com/alexfalkowski/go-service/v2/bytes"
	"github.com/alexfalkowski/go-service/v2/os"
)

// Map is a key-value map.
type Map map[string]string

// Secrets will traverse the map and load any secrets that have been configured.
func (m Map) Secrets(fs *os.FS) error {
	for key, name := range m {
		data, err := fs.ReadSource(name)
		if err != nil {
			return err
		}

		m[key] = bytes.String(data)
	}

	return nil
}
