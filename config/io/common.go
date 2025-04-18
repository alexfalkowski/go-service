package io

import (
	"path/filepath"

	"github.com/alexfalkowski/go-service/env"
	"github.com/alexfalkowski/go-service/os"
	"github.com/alexfalkowski/go-service/strings"
)

// NewCommon for io.
func NewCommon(name env.Name, fs os.FileSystem) *Common {
	location, kind := common(name)

	return &Common{location: location, kind: kind, fs: fs}
}

// Common for io.
type Common struct {
	fs       os.FileSystem
	location string
	kind     string
}

// It will first try to read from the underlying read writer, otherwise it will read a configuration file in commonly defined locations.
func (c *Common) Read() ([]byte, error) {
	if strings.IsEmpty(c.location) {
		return nil, ErrLocationMissing
	}

	return c.fs.ReadFile(c.location)
}

// Write to the underlying read writer.
func (c *Common) Write(data []byte, mode os.FileMode) error {
	if strings.IsEmpty(c.location) {
		return ErrLocationMissing
	}

	return c.fs.WriteFile(c.location, data, mode)
}

// Kind from the underlying read writer, otherwise YAML.
func (c *Common) Kind() string {
	if strings.IsEmpty(c.kind) {
		return "yml"
	}

	return c.kind
}

func common(name env.Name) (string, string) {
	extensions := []string{".yml", ".toml", ".json"}
	for _, extension := range extensions {
		n := name.String()
		file := n + extension
		dirs := []string{
			os.ExecutableDir(),
			filepath.Join(os.UserHomeDir(), ".config", n),
			"/etc/" + n,
		}

		for _, dir := range dirs {
			name := filepath.Join(dir, file)
			if os.PathExists(name) {
				return name, os.PathExtension(name)
			}
		}
	}

	return "", ""
}
