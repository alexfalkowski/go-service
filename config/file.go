package config

import (
	"errors"
	"os"
)

// nolint:gofumpt
const perm = 0600

// ErrMissingConfigFile for config.
var ErrMissingConfigFile = errors.New("missing config file")

// File config location.
func File() (string, error) {
	return FileFromEnv("CONFIG_FILE")
}

// FileFromEnv location.
func FileFromEnv(env string) (string, error) {
	configFile := os.Getenv(env)
	if configFile == "" {
		return "", ErrMissingConfigFile
	}

	return configFile, nil
}

// ReadFile from config location.
func ReadFile() ([]byte, error) {
	configFile, err := File()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(configFile)
}

// WriteFileToEnv location.
func WriteFileToEnv(env string, data []byte) error {
	configFile, err := FileFromEnv(env)
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, perm)
}
