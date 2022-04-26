package config

import (
	"os"
	"path/filepath"
)

// nolint:gofumpt
const perm = 0600

// File config location.
func File() string {
	return FileFromEnv("CONFIG_FILE")
}

// FileFromEnv location.
func FileFromEnv(env string) string {
	return os.Getenv(env)
}

// ReadFile from config location.
func ReadFile() ([]byte, error) {
	return os.ReadFile(filepath.Clean(File()))
}

// WriteFileToEnv location.
func WriteFileToEnv(env string, data []byte) error {
	return os.WriteFile(filepath.Clean(FileFromEnv(env)), data, perm)
}
