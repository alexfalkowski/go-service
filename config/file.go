package config

import (
	"os"
	"path/filepath"
)

const (
	perm       = 0o600
	configfile = "CONFIG_FILE"
)

// File config location.
func File() string {
	return FileFromEnv(configfile)
}

// FileFromEnv location.
func FileFromEnv(env string) string {
	return os.Getenv(env)
}

// ReadFile from config location.
func ReadFile() ([]byte, error) {
	return ReadFileFromEnv(configfile)
}

// ReadFileFromEnv variable of config location.
func ReadFileFromEnv(env string) ([]byte, error) {
	return os.ReadFile(filepath.Clean(FileFromEnv(env)))
}

// WriteFileToEnv location.
func WriteFileToEnv(env string, data []byte) error {
	return os.WriteFile(filepath.Clean(FileFromEnv(env)), data, perm)
}
