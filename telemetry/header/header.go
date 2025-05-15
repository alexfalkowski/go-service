package header

import (
	"github.com/alexfalkowski/go-service/bytes"
	"github.com/alexfalkowski/go-service/os"
)

// Map is a key-value map.
type Map map[string]string

// Secrets will traverse the map and load any secrets that have been configured.
func (m Map) Secrets(fs os.FileSystem) error {
	for key, name := range m {
		if !fs.PathExists(name) {
			continue
		}

		data, err := fs.ReadFile(name)
		if err != nil {
			return err
		}

		m[key] = bytes.String(data)
	}

	return nil
}
