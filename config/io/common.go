package io

import (
	"github.com/alexfalkowski/go-service/v2/env"
	"github.com/alexfalkowski/go-service/v2/os"
	"github.com/alexfalkowski/go-service/v2/strings"
)

// NewCommon for io.
func NewCommon(name env.Name, fs *os.FS) *Common {
	location, kind := common(name, fs)

	return &Common{location: location, kind: kind, fs: fs}
}

// Common for io.
type Common struct {
	fs       *os.FS
	location string
	kind     string
}

// Read from the common location.
func (c *Common) Read() ([]byte, error) {
	if strings.IsEmpty(c.location) {
		return nil, ErrLocationMissing
	}

	return c.fs.ReadFile(c.location)
}

// Kind from the common location.
func (c *Common) Kind() string {
	return c.kind
}

func common(name env.Name, fs *os.FS) (string, string) {
	extensions := []string{".yaml", ".yml", ".toml", ".json"}
	for _, extension := range extensions {
		n := name.String()
		file := n + extension
		dirs := []string{
			fs.ExecutableDir(),
			fs.Join(os.UserHomeDir(), ".config", n),
			"/etc/" + n,
		}

		for _, dir := range dirs {
			name := fs.Join(dir, file)
			if fs.PathExists(name) {
				return name, fs.PathExtension(name)
			}
		}
	}

	return "", ""
}
